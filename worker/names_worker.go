package worker

import (
	"encoding/json"
	"fmt"
	"github.com/Wazzymandias/blockstack-crawler/config"
	"github.com/Wazzymandias/blockstack-crawler/db"
	. "github.com/Wazzymandias/blockstack-crawler/names"
	"github.com/Wazzymandias/blockstack-crawler/routes"
	"github.com/Wazzymandias/blockstack-crawler/storage"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// NameWorker processes name related requests
// By first checking database and storage before
// making remote calls for information.
type NameWorker struct {
	client http.Client

	db      db.BlockstackDB
	storage storage.BlockstackStorage
}

// RetrieveNewNames attempts to compare the latest set of names for each namespace to the set
// found at the provided time. If no previous data exists at the specified time,
// an error is returned. Otherwise a set of new names is returned.
func (rh *NameWorker) RetrieveNewNames(since time.Time) (map[string]map[string]bool, error) {
	var old, current map[string]map[string]bool
	var err error

	if old, err = rh.GetNamesAt(since); err != nil {
		// TODO log error
		return nil, fmt.Errorf("error fetching names for time %s: %+v", since.String(), err)
	}

	if len(old) == 0 {
		// TODO log error
		return nil, fmt.Errorf("user data not found for time %s: cannot compare", since.String())
	}

	current, err = rh.RetrieveNames()

	if err != nil {
		// TODO log err
		return nil, err
	}

	return SelectNew(old, current), nil
}

// GetNames returns the set of all names at current day
func (rh *NameWorker) GetNames() (map[string]map[string]bool, error) {
	return rh.GetNamesAt(time.Now())
}

// GetNamesAt attempts to find and return the set of names for each namespace at the given date
func (rh *NameWorker) GetNamesAt(date time.Time) (result map[string]map[string]bool, err error) {
	result, err = rh.db.GetNamesAt(date)

	if err != nil && err != config.ErrDBKeyNotFound {
		return
	}

	if len(result) > 0 {
		return
	}

	if rh.storage.NamesExistAt(date) {
		return rh.storage.ReadNamesAt(date)
	}

	return nil, nil
}

// RetrieveNames attempts to return the set of names at current date.
// If the names don't exist in database or storage, they are fetch from remote API.
func (rh *NameWorker) RetrieveNames() (result map[string]map[string]bool, err error) {
	result, err = rh.GetNames()

	if err != nil {
		return
	}

	if len(result) > 0 {
		return
	}

	return rh.FetchAndAddNames()
}

// FetchAndAddNames attempts to query the remote API for the set of names for each namespace at
// the current date. If successful, the results are persisted to underlying database and storage.
func (rh *NameWorker) FetchAndAddNames() (names map[string]map[string]bool, err error) {
	names, err = rh.FetchNames()

	if err != nil {
		return
	}

	dbErr, stgErr := rh.AddNames(names)

	if dbErr != nil {
		return nil, fmt.Errorf("error inserting names into database: %+v", dbErr)
	}

	if stgErr != nil {
		return nil, fmt.Errorf("error persisting names to storage: %+v", stgErr)
	}

	return
}

// AddNames persists the set of names for each namespace into the database and storage.
func (rh *NameWorker) AddNames(names map[string]map[string]bool) (dbErr error, stgErr error) {
	return rh.db.PutNames(names), rh.storage.WriteNames(names)
}

// fetchNames processes namespaces concurrently, returning the list
// of names for each namespace and a concatenation of any errors that occur
func (rh *NameWorker) fetchNames(namespaces []string, count int) (<-chan NamespaceNames, error) {
	var err error
	var errs []error

	wg := new(sync.WaitGroup)
	errCh := make(chan error, count)
	namesCh := make(chan NamespaceNames, count)

nsLoop:
	for _, ns := range namespaces {
		select {
		case err = <-errCh:
			errs = append(errs, err)
			break nsLoop
		default:
			wg.Add(1)
			go rh.FetchNamespaceNames(ns, namesCh, errCh, wg)
		}
	}

	wg.Wait()

	close(errCh)
	close(namesCh)

	// If any errors occurred, the error variable would first be set,
	// at which point all subsequent errors would remain in the channel.
	// Any such errors are appended, to be returned as one concatenated error.
	if err != nil {
		for e := range errCh {
			errs = append(errs, e)
		}

		err = fmt.Errorf("one or more errors occurred while fetching names: %+v", errs)
	}

	return namesCh, err
}

func (rh *NameWorker) transformNames(namesCh <-chan NamespaceNames, err error) (map[string]map[string]bool, error) {
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]bool)

	for n := range namesCh {
		result[n.Namespace] = make(map[string]bool)

		for _, name := range n.Names {
			result[n.Namespace][name] = true
		}
	}

	return result, nil
}

func (rh *NameWorker) FetchNames() (map[string]map[string]bool, error) {
	namespaces, err := routes.GetAllNamespaces()

	if err != nil {
		return nil, err
	}

	nsCount := len(namespaces)

	return rh.transformNames(rh.fetchNames(namespaces, nsCount))
}

func (rh *NameWorker) processPages(namespace string, pages <-chan Names, namesCh chan<- NamespaceNames) {
	var nms []string

	for {
		page, valid := <-pages

		if !valid {
			break
		}

		nms = append(nms, page...)
	}

	namesCh <- NamespaceNames{Namespace: namespace, Names: nms}

	close(namesCh)
}

func (rh *NameWorker) fetchNamespaceNames(namespace string,
	pages chan<- Names) []error {

	var errors []error

	u := fmt.Sprintf("%s://%s%s/%s/%s", config.ApiURLScheme, config.ApiHost, routes.GetAllNamespacesPath, namespace, "names")

	wg := new(sync.WaitGroup)
	pagesDone := make(chan struct{}, config.BatchSize)
	errCh := make(chan error, config.BatchSize)

pageLoop:
	for page := uint64(0); ; page += config.BatchSize {
		select {
		case err := <-errCh:
			errors = append(errors, err)
			break pageLoop

		case <-pagesDone:
			break pageLoop

		default:
			for count := uint64(0); count < config.BatchSize; count++ {
				wg.Add(1)
				go rh.processPageRequest(namespace, u, page+count, pages, pagesDone, errCh, wg)
			}
			wg.Wait()
		}
	}

	close(errCh)
	close(pagesDone)

	// close upstream pages channel since all pages must be processed by this point
	close(pages)

	for err := range errCh {
		errors = append(errors, err)
	}

	return errors
}

// fetch asynchronously inserts results to database
func (rh *NameWorker) FetchNamespaceNames(namespace string, out chan<- NamespaceNames, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	pages := make(chan Names, config.BatchSize)
	namesCh := make(chan NamespaceNames, 1)

	go rh.processPages(namespace, pages, namesCh)

	errors := rh.fetchNamespaceNames(namespace, pages)

	if len(errors) > 0 {
		errCh <- fmt.Errorf("error(s) occurred fetching names for namespace %s: %+v", namespace, errors)
		return
	}

	out <- <-namesCh
}

// TODO implement retry
// TODO zap logger
// TODO make sure cleanup even when early returns (need to defer)
func (rh *NameWorker) processPageRequest(namespace string, pageURL string, page uint64,
	pages chan<- Names, done chan<- struct{}, errCh chan<- error, wg *sync.WaitGroup) {

	defer wg.Done()

	req, err := http.NewRequest("GET", pageURL, nil)

	if err != nil {
		errCh <- err
		return
	}

	req.URL.RawQuery = fmt.Sprintf("page=%d", page)

	resp, err := rh.client.Do(req)

	if err != nil {
		errCh <- err
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		errCh <- fmt.Errorf("error reading response body: [err: %v] [body: %v]", err, string(body))
		return
	}

	var pageResults Names
	err = json.Unmarshal(body, &pageResults)

	if err != nil {
		errCh <- fmt.Errorf("error unmarshalling response body: [err: %v] [body: %v]", err, string(body))
		return
	}

	numResults := len(pageResults)

	if numResults == 0 {
		done <- struct{}{}
		return
	}

	pages <- pageResults
}

func (rh *NameWorker) Shutdown() error {
	return rh.db.Shutdown()
}

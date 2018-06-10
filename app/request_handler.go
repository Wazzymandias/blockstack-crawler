package app

import (
	"encoding/json"
	"fmt"
	"github.com/Wazzymandias/blockstack-profile-crawler/config"
	"github.com/Wazzymandias/blockstack-profile-crawler/db"
	"github.com/Wazzymandias/blockstack-profile-crawler/storage"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	getAllNamespaces = "/v1/namespaces"
)

type RequestHandler struct {
	client http.Client

	db      db.BlockstackDB
	storage storage.BlockstackStorage
}

// TODO switch statement that maps request type to url endpoint to hit
func GetAllNamespaces() ([]string, error) {
	var result []string

	url := fmt.Sprintf("%s://%s%s", config.ApiURLScheme, config.ApiURL, getAllNamespaces)
	client := &http.Client{Timeout: config.Timeout}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rh *RequestHandler) RetrieveNewNames(since time.Time) (map[string]map[string]bool, error) {
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

	return selectNewNames(old, current), nil
}

func (rh *RequestHandler) GetNames() (map[string]map[string]bool, error) {
	return rh.GetNamesAt(time.Now())
}

func (rh *RequestHandler) GetNamesAt(date time.Time) (result map[string]map[string]bool, err error) {
	result, err = rh.db.GetNamesAt(date)

	if err != nil {
		return
	}

	if result != nil {
		return
	}

	if rh.storage.NamesExistAt(date) {
		return rh.storage.ReadNamesAt(date)
	}

	return nil, nil
}

func (rh *RequestHandler) RetrieveNames() (result map[string]map[string]bool, err error) {
	result, err = rh.GetNames()

	if err != nil {
		return
	}

	if result != nil {
		return
	}

	return rh.FetchAndAddNames()
}

func (rh *RequestHandler) FetchAndAddNames() (names map[string]map[string]bool, err error) {
	names, err = rh.FetchNames()

	if err != nil {
		return
	}

	if err = rh.AddNames(names); err != nil {
		return
	}

	return
}

func (rh *RequestHandler) AddNames(names map[string]map[string]bool) (err error) {
	if err = rh.db.PutNames(names); err != nil {
		fmt.Println("error occured adding names to db")
		return
	}

	return rh.storage.WriteNames(names)
}

//func (rh *RequestHandler) seed(names map[string]map[string]bool, t time.Time) error {
//	seeding := make(map[string]map[string]bool)
//
//	for k, v := range names {
//		seeding[k] = make(map[string]bool)
//
//		for n := range v {
//			seeding[k][n] = true
//			break
//		}
//	}
//
//	fmt.Println("seed data: ", seeding)
//	fmt.Println(t.AddDate(0, 0, -1))
//	return rh.db.PutNamesAt(names, t.AddDate(0, 0, -1))
//}

func (rh *RequestHandler) fetchNames(namespaces []string, count int) (<-chan names, <-chan error) {
	var errors []error

	wg := new(sync.WaitGroup)
	errCh := make(chan error, count)
	namesCh := make(chan names, count)

nsLoop:
	for _, ns := range namespaces {
		select {
		case err := <-errCh:
			errors = append(errors, err)
			break nsLoop
		default:
			wg.Add(1)
			go rh.FetchNamespaceNames(ns, namesCh, errCh, wg)
		}
	}

	wg.Wait()

	close(errCh)
	close(namesCh)

	return namesCh, errCh
}

func (rh *RequestHandler) transformNames(namesCh <-chan names, errorCh <-chan error) (map[string]map[string]bool, error) {
	var errors []error
	errCount := 0

	for err := range errorCh {
		errors = append(errors, err)
		errCount++
	}

	if errCount > 0 {
		return nil, fmt.Errorf("error(s) occured fetching users: %+v", errors)
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

func (rh *RequestHandler) FetchNames() (map[string]map[string]bool, error) {
	namespaces, err := GetAllNamespaces()

	if err != nil {
		return nil, err
	}

	nsCount := len(namespaces)

	return rh.transformNames(rh.fetchNames(namespaces, nsCount))
}

func (rh *RequestHandler) processPages(namespace string, pages <-chan NamesPage, namesCh chan<- names) {
	var nms []string

	for {
		page, valid := <-pages

		if !valid {
			break
		}

		nms = append(nms, page.UserIDs...)
		fmt.Printf("finished processing page %d for %s\n", page.PageNum, namespace)
	}

	namesCh <- names{Namespace: namespace, Names: nms}

	close(namesCh)
}

func (rh *RequestHandler) fetchNamespaceNames(namespace string,
	pages chan<- NamesPage) []error {

	var errors []error

	u := fmt.Sprintf("%s://%s%s/%s/%s", config.ApiURLScheme, config.ApiURL, getAllNamespaces, namespace, "names")

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
func (rh *RequestHandler) FetchNamespaceNames(namespace string, out chan<- names, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	pages := make(chan NamesPage, config.BatchSize)
	namesCh := make(chan names, 1)

	go rh.processPages(namespace, pages, namesCh)

	errors := rh.fetchNamespaceNames(namespace, pages)

	if len(errors) > 0 {
		errCh <- fmt.Errorf("error(s) occurred fetching names for namespace %s: %+v", namespace, errors)
		return
	}

	out <- <-namesCh
}

type names struct {
	Namespace string
	Names     []string
}

// nsNamesByDate(
// TODO implement retry
// event handler that watches for changes in directory for each namespace name page file input -> db insert output
// event handler with url input and json output
// event handler with json input and file output
// prefix scan over elements to get all names in a namespace
// check in database first - look up key for namespace/date -> namespace/names/date/page.json
// check namespace/names/date directory existence
// check number of files in directory
// if it's not there check file
// if it's not there make a network request
// ns.Names(namespace)
// TODO zap logger
// TODO make sure cleanup even when early returns (need to defer)

func (rh *RequestHandler) processPageRequest(namespace string, pageURL string, page uint64,
	pages chan<- NamesPage, done chan<- struct{}, errCh chan<- error, wg *sync.WaitGroup) {

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

	var pageResults []string
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

	pages <- NamesPage{PageNum: page, UserIDs: pageResults, Count: numResults}
}

func (rh *RequestHandler) Shutdown() error {
	return rh.db.Shutdown()
}
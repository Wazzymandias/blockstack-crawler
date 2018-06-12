// package badger implements database operations for BadgerDB
package badger

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/Wazzymandias/blockstack-crawler/config"
	"github.com/dgraph-io/badger"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
)

// DB is a wrapper with reference to underlying Badger database and
// implements BlockstackDB operations
type DB struct {
	db *badger.DB
}

// New attempts to create a new Badger database instance given path.
func New(path string) (*DB, error) {
	opts := badger.DefaultOptions

	opts.Dir = path
	opts.ValueDir = path

	b, err := badger.Open(opts)

	if err != nil {
		return &DB{}, err
	}

	return &DB{db: b}, nil
}

// Get attempts to return value at given key.
// If an error occurs, it is returned and the result will be
// an empty byte slice.
func (db *DB) Get(key []byte) (result []byte, err error) {
	err = db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)

		if err != nil {
			if err == badger.ErrKeyNotFound {
				// TODO log error
				err = config.ErrDBKeyNotFound
			}

			// TODO log error
			return err
		}

		result, err = item.ValueCopy(result)

		if err != nil {
			// TODO log error
			return err
		}

		return nil
	})

	return
}

// Put attempts to insert a key value pair,
// and returns an error if any occur
func (db *DB) Put(key, value []byte) error {
	txn := db.db.NewTransaction(true)
	defer txn.Discard()

	err := txn.Set(key, value)

	if err != nil {
		return err
	}

	return txn.Commit(nil)
}

// Shutdown closes the database. It is critical
// to call this before closing the application in order
// ensure persistence of all data by flushing to disk.
func (db *DB) Shutdown() error {
	return db.db.Close()
}

// GetNames returns the set of all names at current day
func (db *DB) GetNames() (map[string]map[string]bool, error) {
	return db.GetNamesAt(time.Now())
}

// GetNamesAt returns set of names across all namespaces at a given day.
// The time parameter is used as the key, rounded to start of day in UTC format,
// and converted to Unix epoch time.
// The result is map[NamespaceName][SetOfNames], using a boolean map for
// the name set since Go does not have builtin sets.
// The set of names are stored as gzip compressed bytes, and are uncompressed when
// retrieved from database.
func (db *DB) GetNamesAt(t time.Time) (result map[string]map[string]bool, err error) {
	// time is rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	key := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10))

	value, err := db.Get([]byte(key))

	if err != nil || len(value) == 0 {
		return
	}

	gz, err := gzip.NewReader(ioutil.NopCloser(
		bytes.NewBuffer(value)))

	if err != nil {
		return
	}

	decompressed, err := ioutil.ReadAll(gz)

	if err != nil {
		return
	}

	err = gz.Close()

	if err != nil {
		return
	}

	err = json.Unmarshal(decompressed, &result)

	return
}

// PutNames inserts the set of names for each namespace into the database
// where the key is the current date rounded to start of day.
func (db *DB) PutNames(names map[string]map[string]bool) error {
	return db.PutNamesAt(names, time.Now())
}

// TODO insert the names per namespace as well, otherwise would always need to grab all names
// PutNamesAt attempts to insert the set of names for each namespace into the database using the
// given time as key, rounded to start of day in UTC and converted to Unix epoch time.
// The set of names are stored as gzip compressed bytes, and are uncompressed when
// retrieved from database.
func (db *DB) PutNamesAt(names map[string]map[string]bool, t time.Time) error {
	// time is rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	key := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10))

	value, err := json.Marshal(&names)

	if err != nil {
		return err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(value); err != nil {
		return err
	}

	err = gz.Close()

	if err != nil {
		return err
	}

	return db.Put([]byte(key), value)
}

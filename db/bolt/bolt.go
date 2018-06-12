// package bolt implements operations using Bolt as the underlying database
package bolt

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/coreos/bbolt"
	"io/ioutil"
	"strconv"
	"time"
)

const namesBucket = "names"

// DB is a wrapper with reference to underlying Bolt database and
// implements BlockstackDB operations
type DB struct {
	DB *bolt.DB
}

// New attempts to create a new Bolt database instance given path.
func New(path string) (*DB, error) {
	db, err := bolt.Open(path, 0755, bolt.DefaultOptions)

	if err != nil {
		return &DB{}, err
	}

	return &DB{DB: db}, err
}

// Get attempts to return value at given key.
// If an error occurs, it is returned and the result will be
// an empty byte slice.
func (db *DB) Get(bucket, key []byte) (result []byte, err error) {
	err = db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		if b != nil {
			result = b.Get(key)
		}

		return nil
	})

	return
}

// Put attempts to insert the given key value pair under the
// provided bucket name. The bucket is created if it doesn't
// exist.
func (db *DB) Put(bucket, key, value []byte) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)

		if err != nil {
			return err
		}

		return b.Put(key, value)
	})
}

// Shutdown closes the database and returns an error if any occur.
func (db *DB) Shutdown() error {
	return db.DB.Close()
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

	key := strconv.FormatInt(rounded.Unix(), 10)

	bucket := []byte(namesBucket)

	value, err := db.Get(bucket, []byte(key))

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

// PutNamesAt attempts to insert the set of names for each namespace into the database using the
// given time as key, rounded to start of day in UTC and converted to Unix epoch time.
// The set of names are stored as gzip compressed bytes, and are uncompressed when
// retrieved from database.
func (db *DB) PutNamesAt(names map[string]map[string]bool, t time.Time) error {
	// time is rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	key := strconv.FormatInt(rounded.Unix(), 10)

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

	return db.Put([]byte(namesBucket), []byte(key), buf.Bytes())
}

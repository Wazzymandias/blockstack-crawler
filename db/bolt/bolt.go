package bolt

import (
	"github.com/coreos/bbolt"
	"time"
	"strconv"
	"compress/gzip"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

const namesBucket = "names"

type DB struct {
	db *bolt.DB
}

func New(path string) (*DB, error) {
	db, err := bolt.Open(path, 0755, bolt.DefaultOptions)

	if err != nil {
		return &DB{}, err
	}

	return &DB{db:db}, err
}

func (db *DB) Get(bucket, key []byte) (result []byte, err error) {
	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		if b != nil {
			result = b.Get(key)
		}

		return nil
	})

	return
}

func (db *DB) Put(bucket, key, value []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)

		if err != nil {
			return err
		}

		return b.Put(key, value)
	})
}

func (db *DB) Shutdown() error {
	return db.db.Close()
}

func (db *DB) GetNames() (map[string]map[string]bool, error) {
	return db.GetNamesAt(time.Now())
}

func (db *DB) GetNamesAt(t time.Time) (result map[string]map[string]bool, err error) {
	// time is rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

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

func (db *DB) PutNames(names map[string]map[string]bool) error {
	return db.PutNamesAt(names, time.Now())
}

func (db *DB) PutNamesAt(names map[string]map[string]bool, t time.Time) error {
	// time is rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

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

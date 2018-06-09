package badger

import (
	"github.com/dgraph-io/badger"
	"time"
	"path/filepath"
	"encoding/json"
	"github.com/Wazzymandias/blockstack-profile-crawler/config"
)

type DB struct {
	db *badger.DB
}

func New(path string) (DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = path
	db, err := badger.Open(opts)

	if err != nil {
		return DB{}, err
	}

	return DB{db:db}, nil
}

func (db *DB) Get(key []byte) (result []byte, err error) {
	err = db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)

		if err != nil {
			if err == badger.ErrKeyNotFound {
				// TODO log error
				return nil
			}

			// TODO log error
			return err
		}

		val, err := item.Value()

		if err != nil {
			// TODO log error
			return err
		}

		copy(result, val)

		return nil
	})

	return
}

func (db *DB) Put(key, value []byte) error {
	return nil
}

func (db *DB) GetNames() (map[string]map[string]bool, error) {
	return db.GetNamesAt(time.Now())
}

func (db *DB) GetNamesAt(t time.Time) (result map[string]map[string]bool, err error) {
	// time is rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	key := filepath.Join(config.DataDir, config.NamesDir, string(rounded.Unix()))

	value, err := db.Get([]byte(key))

	if err != nil {
		return
	}

	err = json.Unmarshal(value, &result)

	return
}

func (db *DB) PutNames(names map[string]map[string]bool) error {
	return db.PutNamesAt(names, time.Now())
}

func (db *DB) PutNamesAt(names map[string]map[string]bool, t time.Time) error {
	// time is rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	key := filepath.Join(config.DataDir, config.NamesDir, string(rounded.Unix()))

	value, err := json.Marshal(&names)

	if err != nil {
		return err
	}

	return db.Put([]byte(key), value)
}
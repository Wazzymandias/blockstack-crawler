package badger

import "github.com/dgraph-io/badger"

type DB struct {
	db *badger.DB
}

func (db *DB) Get(key interface{}) (interface{}, bool) {
	return nil, false
}

func (db *DB) Put(key, value interface{}) {

}

package db

type DB interface {
	Put(key, value interface{})
	Get(key interface{}) (interface{}, bool)
}

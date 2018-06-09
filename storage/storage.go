package storage

import "time"

type Storage interface {
	Exists(key interface{}) bool
	Find(path string) (interface{}, bool)
}

type BlockstackStorage interface {
	WriteNames(names map[string]map[string]bool) error
	WriteNamesAt(names map[string]map[string]bool, time time.Time) error
	ReadNames() error
	ReadNamesAt(time time.Time) (map[string]map[string]bool, error)
}

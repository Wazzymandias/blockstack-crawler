package storage

import "time"

// Storage defines basic operations
type Storage interface {
	Exists(path string) bool
	Find(path string) ([]byte, error)
	Write(path string, value []byte) error
	Read(path string) ([]byte, error)
}

type BlockstackStorage interface {
	Storage

	WriteNames(names map[string]map[string]bool) error
	WriteNamesAt(names map[string]map[string]bool, t time.Time) error
	ReadNames() (map[string]map[string]bool, error)
	ReadNamesAt(t time.Time) (map[string]map[string]bool, error)
	NamesExist() bool
	NamesExistAt(t time.Time) bool
}

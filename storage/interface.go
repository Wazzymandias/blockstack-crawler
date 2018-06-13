package storage

import "time"

// Storage defines basic operations that need to be implemented by storage drivers
type Storage interface {
	Exists(path string) bool
	Find(path string) ([]byte, error)
	Write(path string, value []byte) error
	Read(path string) ([]byte, error)
}

// BlockstackStorage defines Blockstack related operations that need to be implemented by storage drivers
type BlockstackStorage interface {
	Storage

	WriteNames(names map[string][]string) error
	WriteNamesAt(names map[string][]string, t time.Time) error
	ReadNames() (map[string][]string, error)
	ReadNamesAt(t time.Time) (map[string][]string, error)
	NamesExist() bool
	NamesExistAt(t time.Time) bool
}

package db

import "time"

type DB interface {
	Shutdown() error
}

type BlockstackDB interface {
	DB

	GetNames() (map[string]map[string]bool, error)
	GetNamesAt(t time.Time) (map[string]map[string]bool, error)
	PutNames(names map[string]map[string]bool) error
	PutNamesAt(names map[string]map[string]bool, t time.Time) error
}

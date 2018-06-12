package db

import "time"

// DB interface enforces Shutdown method, which may be
// required for certain databases to ensure persistence of all operations
// (e.g. flushing to disk from memory), and
// is useful if signal is sent to application to gracefully terminate
type DB interface {
	Shutdown() error
}

// BlockstackDB provides interface for Blockstack related operations
type BlockstackDB interface {
	DB

	GetNames() (map[string]map[string]bool, error)
	GetNamesAt(t time.Time) (map[string]map[string]bool, error)
	PutNames(names map[string]map[string]bool) error
	PutNamesAt(names map[string]map[string]bool, t time.Time) error
}

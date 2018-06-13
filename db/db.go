// Package db implements interfaces to different database types and provides specification of
// BlockstackDB and DB operations, such that any new database can be used as long as it implements those methods.
package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Wazzymandias/blockstack-crawler/config"
	"github.com/Wazzymandias/blockstack-crawler/db/badger"
	"github.com/Wazzymandias/blockstack-crawler/db/bolt"
)

// NewBadgerDB instantiates and returns a Badger DB instance based on the
// provided path. If one or more directories don't exist,
// it will attempt to create them, and returns an error
// if any occur.
func NewBadgerDB(path string) (BlockstackDB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return nil, fmt.Errorf("unable to create new badger DB, cannot create path %s: %+v", path, err)
		}
	}

	return badger.New(path)
}

// NewBoltDB instantiates and returns a Bolt DB instance based on the
// provided path. If one or more directories don't exist,
// it will attempt to create them, and returns an error
// if any occur.
func NewBoltDB(path string) (BlockstackDB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return nil, fmt.Errorf("unable to create new bolt DB, cannot create path %s: %+v", path, err)
		}
	}

	return bolt.New(path)
}

// New attempts to instantiate a new database instance using the configuration
// variables provided.
func New() (BlockstackDB, error) {
	switch config.DBFromString(config.DatabaseType) {
	case config.Bolt:
		return NewBoltDB(filepath.Join(config.DataDir, config.DBDir, config.BoltDBFile))
	case config.Badger:
		return NewBadgerDB(filepath.Join(config.DataDir, config.DBDir))
	default:
		return nil, fmt.Errorf("unsupported database type entered [%v]", config.DatabaseType)
	}
}

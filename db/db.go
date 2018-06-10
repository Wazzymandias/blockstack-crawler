package db

import (
	"github.com/Wazzymandias/blockstack-profile-crawler/config"
	"fmt"
	"os"
	"path/filepath"
	"github.com/Wazzymandias/blockstack-profile-crawler/db/bolt"
)

func NewBoltDB(path string) (BlockstackDB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("unable to create new badger DB, cannot create path %s: %+v", path, err)
		}
	}

	return bolt.New(path)
}

func New() (BlockstackDB, error) {
	switch config.DatabaseType {
	case config.Bolt:
		return NewBoltDB(filepath.Join(config.DataDir, config.DBDir, "bolt.db"))
	default:
		return nil, fmt.Errorf("unsupported database type entered [%v]", config.DatabaseType)
	}
}
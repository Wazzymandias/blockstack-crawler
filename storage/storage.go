// Package storage provides interfaces to different storage types to persist data
package storage

import (
	"fmt"

	"github.com/Wazzymandias/blockstack-crawler/config"
)

// New returns a reference to a storage driver that implements blockstack storage operations.
// If the storage type specified in configuration variables is invalid, or there is an error
// instantiating the storage driver, an error is returned.
func New() (BlockstackStorage, error) {
	switch config.StorageFromString(config.StorageType) {
	case config.Local:
		return NewLocal(config.DataDir, 0644, 0755)
	default:
		return nil, fmt.Errorf("unsupported storage type [%v]", config.StorageType)
	}
}

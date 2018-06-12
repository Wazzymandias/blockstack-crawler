// package storage provides interfaces to different storage types to persist data
package storage

import (
	"fmt"

	"github.com/Wazzymandias/blockstack-crawler/config"
)

func New() (BlockstackStorage, error) {
	switch config.StorageFromString(config.StorageType) {
	case config.Local:
		return NewLocal(config.DataDir)
	default:
		return nil, fmt.Errorf("unsupported storage type [%v]", config.StorageType)
	}
}

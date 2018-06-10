package storage

import (
	"fmt"

	"github.com/Wazzymandias/blockstack-profile-crawler/config"
)

func New() (BlockstackStorage, error) {
	switch config.StorageType {
	case config.Local:
		return NewLocal(config.DataDir)
	default:
		return nil, fmt.Errorf("unsupported storage type [%v]", config.StorageType)
	}
}

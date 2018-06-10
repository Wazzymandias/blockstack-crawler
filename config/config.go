package config

import (
	"fmt"
	"os"
	"time"
)

const (
	ProgramName = "bpc"

	DefaultApiURL       = "core.blockstack.org"
	DefaultApiURLScheme = "https"

	DefaultTimeout   = 120 * time.Second
	DefaultBatchSize = 50

	DefaultDBType      = Badger
	DefaultDBDir       = "db"
	DefaultStorageType = Local

	NamesDir  = "names"
	NamesJSON = "names.json"
)

var (
	DefaultDataDir = os.ExpandEnv(fmt.Sprintf("$HOME/.%s/data", ProgramName))

	ApiURL              = DefaultApiURL
	Timeout             = DefaultTimeout
	ApiURLScheme        = DefaultApiURLScheme
	BatchSize    uint64 = DefaultBatchSize

	DataDir      = DefaultDataDir
	DatabaseType = DefaultDBType
	StorageType  = DefaultStorageType
	DBDir        = DefaultDBDir
)

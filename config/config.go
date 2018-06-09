package config

import (
	"time"
)

const (
	defaultApiURL       = "core.blockstack.org"
	defaultTimeout      = 120 * time.Second
	defaultApiURLScheme = "https"
	defaultBatchSize    = 50
	ChannelCap          = 250

	NamesDir = "names"

	StorageFileType = "json"

	defaultDBType = Badger
)

var (
	ApiURL              = defaultApiURL
	Timeout             = defaultTimeout
	ApiURLScheme        = defaultApiURLScheme
	BatchSize    uint64 = defaultBatchSize

	DataDir string
	DBType  = defaultDBType
)

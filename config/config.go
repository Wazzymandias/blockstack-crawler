package config

import (
	"github.com/Wazzymandias/blockstack-profile-crawler/db"
	"time"
)

const (
	defaultApiURL       = "core.blockstack.org"
	defaultTimeout      = 120 * time.Second
	defaultApiURLScheme = "https"
	defaultBatchSize    = 50
	ChannelCap          = 250

	UsersDir        = "users"
	StorageFileType = "json"

	defaultDBType = db.Badger
)

var (
	ApiURL              = defaultApiURL
	Timeout             = defaultTimeout
	ApiURLScheme        = defaultApiURLScheme
	BatchSize    uint64 = defaultBatchSize

	DataDir string
	DBType  = defaultDBType
)

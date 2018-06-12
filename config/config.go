// Package config stores all configuration related variables and constants, and
// stores any arguments passed by the user
package config

import (
	"fmt"
	"os"
	"time"
)

const (
	// ProgramName is name of program, references in CLI when displaying commands
	ProgramName = "blockstack-crawler"

	// DefaultAPIHost
	DefaultAPIHost   = "core.blockstack.org"
	DefaultAPIScheme = "https"

	DefaultTimeout   = 120 * time.Second
	DefaultBatchSize = 50

	DefaultAPIPort = 80

	DefaultDBType      = "bolt"
	DefaultDBDir       = "db"
	DefaultStorageType = "local"

	DefaultOutputFormat = "json"

	NamesDir  = "names"
	NamesJSON = "names.json"

	BoltDBFile = "bolt.db"
)

var (
	// DefaultDataDir
	DefaultDataDir = os.ExpandEnv(fmt.Sprintf("$HOME/.%s/data", ProgramName))

	ApiHost      = DefaultAPIHost
	Timeout      = DefaultTimeout
	ApiURLScheme = DefaultAPIScheme
	ApiPort      = DefaultAPIPort

	BatchSize    uint64 = DefaultBatchSize

	DataDir      = DefaultDataDir
	DatabaseType string
	StorageType  string
	DBDir        = DefaultDBDir

	NewUsersSince string
	OutputFormat = DefaultOutputFormat

	OutputFile string
)

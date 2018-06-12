package config

import "errors"

type DBEngine int

const (
	Bolt DBEngine = iota
	Badger
	InvalidDB
)

func DBFromString(s string) DBEngine {
	switch s {
	case "badger", "Badger", "BADGER":
		return Badger
	case "bolt", "Bolt", "BOLT":
		return Bolt
	default:
		return InvalidDB
	}
}

var ErrDBKeyNotFound = errors.New("key not found")

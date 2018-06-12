package config

type Store int

const (
	Local Store = iota
	InvalidStore
)

func StorageFromString(s string) Store {
	switch s {
	case "local", "Local", "LOCAL":
		return Local
	default:
		return InvalidStore
	}
}

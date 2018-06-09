package db

import "github.com/Wazzymandias/blockstack-profile-crawler/db/badger"

func NewBadgerDB(path string) (badger.DB, error) {
	return badger.New(path)
}

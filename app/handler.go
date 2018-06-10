package app

import (
	"github.com/Wazzymandias/blockstack-profile-crawler/storage"
	"github.com/Wazzymandias/blockstack-profile-crawler/db"
)

func NewRequestHandler() (*RequestHandler, error) {

	store, err := storage.New()

	if err != nil {
		return nil, err
	}

	rhDB, err := db.New()

	if err != nil {
		return nil, err
	}

	return &RequestHandler{storage: store, db: rhDB}, nil
}

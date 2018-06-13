// Package worker implements various workers to process requests sent by the user
package worker

import (
	"github.com/Wazzymandias/blockstack-crawler/config"
	"github.com/Wazzymandias/blockstack-crawler/db"
	"github.com/Wazzymandias/blockstack-crawler/storage"
	"net/http"
)

// NewNameWorker initializes a worker with new storage and
// database based on configuration variables specified.
// The worker can perform name-related operations using Blockstack API
// and persist any relevant data into database and storage.
func NewNameWorker() (*NameWorker, error) {
	store, err := storage.New()

	if err != nil {
		return &NameWorker{}, err
	}

	d, err := db.New()

	if err != nil {
		return &NameWorker{}, err
	}

	cl := http.Client{Timeout: config.Timeout}

	return &NameWorker{storage: store, db: d, client: cl}, nil
}

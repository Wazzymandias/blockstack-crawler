// worker implements various workers to process requests sent by the user
package worker

import (
	"github.com/Wazzymandias/blockstack-crawler/storage"
	"github.com/Wazzymandias/blockstack-crawler/db"
	"net/http"
	"github.com/Wazzymandias/blockstack-crawler/config"
)

// NewNameWorker initializes a worker with new storage and
// database based on configuration variables specified.
// The worker can perform name-related operations using Blockstack API
// and persist any relevant data into database and storage.
func NewNameWorker() (*NamesWorker, error) {
	store, err := storage.New()

	if err != nil {
		return &NamesWorker{}, err
	}

	d, err := db.New()

	if err != nil {
		return &NamesWorker{}, err
	}

	cl := http.Client{Timeout:config.Timeout}

	return &NamesWorker{storage:store, db:d, client:cl}, nil
}
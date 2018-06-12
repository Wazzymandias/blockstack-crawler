package worker

import (
	"github.com/Wazzymandias/blockstack-crawler/config"
	"testing"
)

func TestNameWorker(t *testing.T) {
	config.DataDir = "/tmp"

	_, err := NewNameWorker()

	if err != nil {
		t.Fatal(err)
	}
}

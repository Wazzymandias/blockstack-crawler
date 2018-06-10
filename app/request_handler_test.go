package app

import (
	"testing"
	"github.com/Wazzymandias/blockstack-profile-crawler/config"
	"fmt"
)

func TestRequestHandler(t *testing.T) {
	config.DataDir = "/tmp"

	rh, err := NewRequestHandler()

	if err != nil {
		t.Fatal(err)
	}

	err = rh.db.Put([]byte("foo"), []byte("bar"))

	if err != nil {
		t.Fatal(err)
	}

	v, err := rh.db.Get([]byte("foo"))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("get while db is up: ", string(v))

	err = rh.Shutdown()

	if err != nil {
		t.Fatal(err)
	}

	rh, err = NewRequestHandler()

	if err != nil {
		t.Fatal(err)
	}

	v, err = rh.db.Get([]byte("foo"))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("get after db shut down: ", string(v))
}

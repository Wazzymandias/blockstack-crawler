package db

import (
	"fmt"
	"github.com/Wazzymandias/blockstack-crawler/config"
	"path/filepath"
	"testing"
	"time"
)

func TestNewBoltDB(t *testing.T) {
	//config.DataDir = "/tmp"
	b, err := NewBoltDB(filepath.Join(config.DataDir, config.DBDir, "bolt.db"))

	if err != nil {
		t.Fatal(err)
	}

	v, err := b.GetNamesAt(time.Now().AddDate(0, 0, -2))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(time.Now().AddDate(0, 0, -2))

	if len(v) > 0 {
		fmt.Println("found v:", len(v))
	}

	layout := "2006-01-02"

	ti, err := time.Parse(layout, "2018-06-08")

	fmt.Println("time:", ti)

	v2, err := b.GetNamesAt(ti)

	if len(v2) > 0 {
		fmt.Println("found v2:", len(v2))
	}
}

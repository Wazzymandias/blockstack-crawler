package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Wazzymandias/blockstack-profile-crawler/config"
)

type Local struct {
	perms os.FileMode
}

// TODO validate directory exists and has read/write permissions
func NewLocal(dataDir string) (*Local, error) {
	return &Local{perms: 0644}, nil
}

func (l *Local) Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (l *Local) Find(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	if !l.Exists(path) {
		return nil, fmt.Errorf("path %s does not exist", path)
	}

	return l.Read(path)
}

// TODO pass file permissions (maybe make that specific to local storage though)
func (l *Local) Write(path string, value []byte) error {
	return ioutil.WriteFile(path, value, 0644)
}

func (l *Local) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (l *Local) WriteNames(names map[string]map[string]bool) error {
	return l.WriteNamesAt(names, time.Now())
}

// TODO attempt to create all directories if they don't exist
func (l *Local) WriteNamesAt(names map[string]map[string]bool, t time.Time) error {
	// rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	// dataDir + names + time + names.json
	path := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10), config.NamesJSON)
	dir := filepath.Dir(path)

	if !l.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf(
				"one or more directories did not exist and could not be created [%s]: %+v", path, err)
		}
	}

	nmBytes, err := json.Marshal(&names)

	if err != nil {
		return fmt.Errorf("error marshalling names to json: %+v", err)
	}

	return ioutil.WriteFile(path, nmBytes, 0644)
}

// ReadNames attempts to read all names at current time, where
// current time is today's date rounded to start of day and converted
// to Unix Epoch time decimal string
func (l *Local) ReadNames() (map[string]map[string]bool, error) {
	return l.ReadNamesAt(time.Now())
}

// ReadNamesAt attempts to read all names at given time, where
// the time parameter is rounded to start of day and converted
// to Unix Epoch time decimal string
func (l *Local) ReadNamesAt(t time.Time) (map[string]map[string]bool, error) {
	// rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	path := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10), config.NamesJSON)

	rdBytes, err := l.Read(path)

	if err != nil {
		return nil, err
	}

	var result map[string]map[string]bool

	err = json.Unmarshal(rdBytes, &result)

	return result, err
}

func (l *Local) NamesExist() bool {
	return l.NamesExistAt(time.Now())
}

func (l *Local) NamesExistAt(t time.Time) bool {
	// rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	path := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10), config.NamesJSON)

	return l.Exists(path)
}

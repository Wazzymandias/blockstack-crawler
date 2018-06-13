package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Wazzymandias/blockstack-crawler/config"
)

// Local is a wrapper struct that implements storage and blockstack storage
// related operations, with file and directory permissions used
// when writing to storage
type Local struct {
	filePerms os.FileMode
	dirPerms os.FileMode
}

// NewLocal returns struct wrapper to operate on local storage
func NewLocal(dataDir string, filePerms os.FileMode, dirPerms os.FileMode) (*Local, error) {
	// TODO validate directory exists and has read/write permissions
	return &Local{filePerms: filePerms, dirPerms:dirPerms}, nil
}

// Exists checks local storage for existence of path
func (l *Local) Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Find attempts to lookup and return the data at specified path.
// If the path is invalid or doesn't exist, an error is returned.
func (l *Local) Find(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	if !l.Exists(path) {
		return nil, fmt.Errorf("path %s does not exist", path)
	}

	return l.Read(path)
}

// Write will write data to the path specified and returns an error if any occur.
func (l *Local) Write(path string, value []byte) error {
	return ioutil.WriteFile(path, value, l.filePerms)
}

// Read will read data from the path specified and returns an error if any occur.
func (l *Local) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// WriteNames will write the set of names to the directory of the current date,
// which is a rounded-to-day Unix epoch time
func (l *Local) WriteNames(names map[string][]string) error {
	return l.WriteNamesAt(names, time.Now())
}

// WriteNamesAt will write the set of names to the directory
func (l *Local) WriteNamesAt(ns map[string][]string, t time.Time) error {
	// rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	// dataDir + names + time + names.json
	path := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10), config.NamesJSON)
	dir := filepath.Dir(path)

	if !l.Exists(dir) {
		if err := os.MkdirAll(dir, l.dirPerms); err != nil {
			return fmt.Errorf(
				"one or more directories did not exist and could not be created [%s]: %+v", path, err)
		}
	}

	nmBytes, err := json.Marshal(&ns)

	if err != nil {
		return fmt.Errorf("error marshalling names to json: %+v", err)
	}

	return ioutil.WriteFile(path, nmBytes, l.filePerms)
}

// ReadNames attempts to read all names at current time, where
// current time is today's date rounded to start of day and converted
// to Unix Epoch time decimal string
func (l *Local) ReadNames() (map[string][]string, error) {
	return l.ReadNamesAt(time.Now())
}

// ReadNamesAt attempts to read all names at given time, where
// the time parameter is rounded to start of day and converted
// to Unix Epoch time decimal string
func (l *Local) ReadNamesAt(t time.Time) (map[string][]string, error) {
	// rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	path := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10), config.NamesJSON)

	rdBytes, err := l.Read(path)

	if err != nil {
		return nil, err
	}

	var result map[string][]string

	err = json.Unmarshal(rdBytes, &result)

	return result, err
}

// NamesExist checks if names exist at the current date
func (l *Local) NamesExist() bool {
	return l.NamesExistAt(time.Now())
}

// NamesExistAt checks if names exist at the specified time, rounded to start of date in UTC
func (l *Local) NamesExistAt(t time.Time) bool {
	// rounded to start of day
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	path := filepath.Join(config.DataDir, config.NamesDir, strconv.FormatInt(rounded.Unix(), 10), config.NamesJSON)

	return l.Exists(path)
}

// Package routes contains Blockstack API related paths and functions to make HTTP requests
// using specified API configuration settings.
package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Wazzymandias/blockstack-crawler/config"
	"io/ioutil"
	"net/http"
)

const (
	GetAllNamespacesPath = "/v1/namespaces"
)

// TODO switch statement that maps request type to url endpoint to hit
func GetAllNamespaces() ([]string, error) {
	var result []string

	url := fmt.Sprintf("%s://%s%s", config.ApiURLScheme, config.ApiHost, GetAllNamespacesPath)
	client := &http.Client{Timeout: config.Timeout}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

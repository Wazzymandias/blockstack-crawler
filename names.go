package main

import (
	"encoding/json"
	"fmt"
	"github.com/Wazzymandias/blockstack-profile-crawler/config"
	"io/ioutil"
	"net/http"
)

// GetAll returns list of all names
func GetAll() error {
	// pages start at 0 for blockstack
	page := 0

	for {
		names, err := getNamesFromPage(page)

		if err != nil {
			return err
		}

		numNames := len(names)

		if numNames == 0 {
			break
		}

		page++
	}

	return nil
}

func getNamesFromPage(page int) ([]string, error) {
	var pageNames []string

	url := fmt.Sprintf("%s/%s", config.ApiURL, "v1/names")

	client := &http.Client{Timeout: config.Timeout}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Form.Set("page", string(page))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, pageNames)

	if err != nil {
		return nil, err
	}

	return pageNames, nil
}

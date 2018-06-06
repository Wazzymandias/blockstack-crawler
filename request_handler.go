package main

import (
	"fmt"
	"github.com/Wazzymandias/blockstack-profile-crawler.go/config"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"
)

const (
	getAllNamespaces = "/v1/namespaces"
)

type RequestHandler struct {
	client http.Client
}

func GetAllNamespaces() ([]string, error) {
	var result []string

	url := fmt.Sprintf("%s://%s%s", config.ApiURLScheme, config.ApiURL, getAllNamespaces)
	client := &http.Client{Timeout:config.Timeout}

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

func GetNamespaceNames(namespace string) ([]string, error) {
	fmt.Println("namespace: ", namespace)

	var result []string

	u := fmt.Sprintf("%s://%s%s/%s/%s", config.ApiURLScheme, config.ApiURL, getAllNamespaces, namespace, "names")
	client := &http.Client{Timeout:config.Timeout}

	page := 0

	for {
		fmt.Println("page: ", page)
		
		var pageNames []string

		data := url.Values{}
		data.Set("page", string(page))

		req, err := http.NewRequest("GET", u, strings.NewReader(data.Encode()))

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

		err = json.Unmarshal(body, &pageNames)

		if err != nil {
			return nil, err
		}

		numNames := len(pageNames)

		if numNames == 0 {
			break
		}

		result = append(result, pageNames...)

		page++
	}

	return result, nil
}
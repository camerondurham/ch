package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	Repository = "camerondurham/ch"
	ApiPath    = "https://api.github.com/repos/%s/releases/latest"
)

func LatestVersion(repository string) string {
	return fmt.Sprintf(ApiPath, repository)
}

func GetLatestVersion() (string, error) {

	url := LatestVersion(Repository)
	data, err := GetRequest(url)
	if err != nil {
		DebugPrint(fmt.Sprintf("error making GET request: %v", err))
		return "", err
	}

	tagName := data["tag_name"].(string)
	return tagName, nil
}

func GetRequest(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		DebugPrint(fmt.Sprintf("error making API request for latest release: %v", err))
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		DebugPrint(fmt.Sprintf("error reading API response body: %v", err))
		return nil, err
	}
	var data map[string]interface{}

	if err = json.Unmarshal(body, &data); err != nil {
		DebugPrint(fmt.Sprintf("error parsing json: %v", err))
		return nil, err
	}

	return data, nil
}

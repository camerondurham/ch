package util

import (
	"encoding/json"
	"fmt"
	"github.com/camerondurham/ch/version"
	"io"
	"net/http"
)

const (
	RepositoryUrl = "https://github.com/camerondurham/ch"
	Repository    = "camerondurham/ch"
	ApiPath       = "https://api.github.com/repos/%s/releases/latest"
)

type callback func(string) (map[string]interface{}, error)

func LatestVersion(repository string) string {
	return fmt.Sprintf(ApiPath, repository)
}

func GetLatestVersion(getRequest callback) (string, error) {

	url := LatestVersion(Repository)
	data, err := getRequest(url)
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

func CheckLatestVersion() {
	latestVersion, err := GetLatestVersion(GetRequest)
	if err != nil {
		DebugPrint(fmt.Sprintf("ignoring version check since error occured when retrieving latest version: %v\n", err))
	} else if version.PkgVersion != "" && latestVersion != version.PkgVersion {
		fmt.Printf(
			"\tYou are running version %s but the latest version is %s.\n"+
				"\tRun `ch upgrade` for upgrade instructions.\n",
			version.PkgVersion,
			latestVersion)
	} else {
		DebugPrint(fmt.Sprintf("local package version: %s\nlatest version: %s\n", version.PkgVersion, latestVersion))
	}
}

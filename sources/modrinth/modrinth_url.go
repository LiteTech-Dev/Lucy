package modrinth

import (
	"encoding/json"
	"io"
	"lucy/apitypes"
	"lucy/syntaxtypes"
	"net/http"
	"net/url"
)

func constructProjectVersionsUrl(slug syntaxtypes.PackageName) (urlString string) {
	urlString, _ = url.JoinPath(
		"https://api.modrinth.com/v2/project",
		string(slug),
		"version",
	)
	return
}

// TODO: Refactor ConstructProjectUrl() to private function

func ConstructProjectUrl(packageName syntaxtypes.PackageName) (url string) {
	return "https://api.modrinth.com/v2/project/" + string(packageName)
}

func GetProjectByName(packageName syntaxtypes.PackageName) (
	project *apitypes.ModrinthProject,
	err error,
) {
	res, err := http.Get(ConstructProjectUrl(packageName))
	if err != nil {
		return
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	project = &apitypes.ModrinthProject{}
	err = json.Unmarshal(data, project)
	return
}

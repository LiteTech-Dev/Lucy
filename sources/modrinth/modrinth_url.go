package modrinth

import (
	"encoding/json"
	"io"
	"lucy/apitypes"
	"lucy/lucytypes"
	"net/http"
	"net/url"
)

func constructProjectVersionsUrl(slug lucytypes.PackageName) (urlString string) {
	urlString, _ = url.JoinPath(
		"https://api.modrinth.com/v2/project",
		string(slug),
		"version",
	)
	return
}

func constructProjectUrl(packageName lucytypes.PackageName) (url string) {
	return "https://api.modrinth.com/v2/project/" + string(packageName)
}

func GetProjectByName(packageName lucytypes.PackageName) (
	project *apitypes.ModrinthProject,
	err error,
) {
	res, err := http.Get(constructProjectUrl(packageName))
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

package modrinth

import (
	"lucy/syntaxtypes"
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

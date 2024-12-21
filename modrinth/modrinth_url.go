package modrinth

import (
	"lucy/syntax"
	"net/url"
)

func constructProjectVersionsUrl(slug syntax.PackageName) (urlString string) {
	urlString, _ = url.JoinPath(
		"https://api.modrinth.com/v2/project",
		string(slug),
		"version",
	)
	return
}

// TODO: Refactor ConstructProjectUrl() to private function

func ConstructProjectUrl(packageName syntax.PackageName) (url string) {
	return "https://api.modrinth.com/v2/project/" + string(packageName)
}

package modrinth

import "net/url"

func constructProjectVersionsUrl(slug string) (urlString string) {
	urlString, _ = url.JoinPath(
		"https://api.modrinth.com/v2/project",
		slug,
		"version",
	)
	return
}

// TODO: Refactor ConstructProjectUrl() to private function

func ConstructProjectUrl(packageName string) (url string) {
	return "https://api.modrinth.com/v2/project/" + packageName
}

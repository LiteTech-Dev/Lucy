package modrinth

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"lucy/probe"
	"lucy/types"
	"net/http"
	url2 "net/url"
	"reflect"
	"strings"
)

// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/

func GetNewestProjectVersion(slug string) (newestVersion types.ModrinthProjectVersion) {
	newestVersion = types.ModrinthProjectVersion{}
	versions := GetProjectVersions(slug)
	serverInfo := probe.GetServerInfo()
	for _, version := range versions {
		for _, gameVersion := range version.GameVersions {
			if gameVersion == serverInfo.Executable.GameVersion &&
				version.VersionType == "release" &&
				version.DatePublished.After(newestVersion.DatePublished) {
				newestVersion = version
			}
		}
	}
	if reflect.DeepEqual(newestVersion, types.ModrinthProjectVersion{}) {
		errors.New("no available version found")
	}

	return
}

func GetProjectVersions(slug string) (versions []types.ModrinthProjectVersion) {
	res, _ := http.Get(constructProjectVersionsUrl(slug))
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &versions)
	return
}

func GetProjectId(slug string) (id string) {
	res, _ := http.Get(ConstructProjectUrl(slug))
	modrinthProject := types.ModrinthProject{}
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &modrinthProject)
	id = modrinthProject.Id
	return
}

// TODO: Search() is way too long, refactor

func Search(
	platform string,
	packageName string,
	showClientPackage bool,
	indexBy string,
) (result types.ModrinthSearchResults) {
	// Construct the search url
	const (
		facetsCategoryAll    = `["categories:'forge'","categories:'fabric'","categories:'quilt'","categories:'liteloader'","categories:'modloader'","categories:'rift'","categories:'neoforge'"]`
		facetsCategoryForge  = `["categories:'forge'"]`
		facetsCategoryFabric = `["categories:'fabric'"]`
		facetsServerOnly     = `["server_side:optional","server_side:required"],["client_side:optional","client_side:required","client_side:unsupported"]`
		facetsShowClient     = `["server_side:optional","server_side:required","server_side:unsupported"],["client_side:optional","client_side:required","client_side:unsupported"]`
		facetsTypeMod        = `["project_type:'mod'"]`
		urlTemplate          = `https://api.modrinth.com/v2/search?query={{.packageName}}&limit=100&index={{.indexBy}}&facets={{.facetsEncoded}}`
	)
	var facetsArray []string
	switch platform {
	case "all":
		facetsArray = append(facetsArray, facetsCategoryAll)
	case "forge":
		facetsArray = append(facetsArray, facetsCategoryForge)
	case "fabric":
		facetsArray = append(facetsArray, facetsCategoryFabric)
	}
	if !showClientPackage {
		facetsArray = append(facetsArray, facetsServerOnly)
	} else {
		facetsArray = append(facetsArray, facetsShowClient)
	}
	facetsArray = append(facetsArray, facetsTypeMod)
	facetsEncoded := url2.QueryEscape(
		"[" + strings.Join(
			facetsArray,
			",",
		) + "]",
	)

	templateUrl, _ := template.New("template_url").Parse(urlTemplate)
	urlBuilder := strings.Builder{}
	_ = templateUrl.Execute(
		&urlBuilder,
		map[string]string{
			"packageName":   packageName,
			"indexBy":       indexBy,
			"facetsEncoded": facetsEncoded,
		},
	)

	// Make the call to Modrinth API
	resp, _ := http.Get(urlBuilder.String())
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	resp.Body.Close()

	return
}

package modrinth

import (
	"encoding/json"
	"html/template"
	"io"
	"lucy/probe"
	"lucy/syntax"
	"lucy/types"
	"net/http"
	url2 "net/url"
	"strings"
)

// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/

func GetNewestProjectVersion(slug syntax.PackageName) (newestVersion *types.ModrinthProjectVersion) {
	newestVersion = nil
	versions := getProjectVersions(slug)
	serverInfo := probe.GetServerInfo()
	for _, version := range versions {
		for _, gameVersion := range version.GameVersions {
			if gameVersion == serverInfo.Executable.GameVersion &&
				version.VersionType == "release" &&
				(newestVersion == nil || version.DatePublished.After(newestVersion.DatePublished)) {
				newestVersion = version
			}
		}
	}

	if newestVersion == nil {
		println("No suitable version found for", slug)
	}

	return
}

func getProjectVersions(slug syntax.PackageName) (versions []*types.ModrinthProjectVersion) {
	res, _ := http.Get(constructProjectVersionsUrl(slug))
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &versions)
	return
}

func GetProjectId(slug syntax.PackageName) (id string) {
	res, _ := http.Get(ConstructProjectUrl(slug))
	modrinthProject := types.ModrinthProject{}
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &modrinthProject)
	id = modrinthProject.Id
	return
}

// TODO: Search() is way too long, refactor

func Search(
	platform syntax.Platform,
	packageName syntax.PackageName,
	showClientPackage bool,
	indexBy string,
) (result *types.ModrinthSearchResults) {
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
	case syntax.AllPlatform:
		facetsArray = append(facetsArray, facetsCategoryAll)
	case syntax.Forge:
		facetsArray = append(facetsArray, facetsCategoryForge)
	case syntax.Fabric:
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
			"packageName":   string(packageName),
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

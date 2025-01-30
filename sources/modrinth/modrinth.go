package modrinth

import (
	"encoding/json"
	"html/template"
	"io"
	"lucy/apitypes"
	"lucy/lucytypes"
	"lucy/probe"
	"lucy/syntaxtypes"
	"net/http"
	url2 "net/url"
	"strings"
)

// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/

func GetNewestProjectVersion(slug syntaxtypes.PackageName) (newestVersion *apitypes.ModrinthProjectVersion) {
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

func getProjectVersions(slug syntaxtypes.PackageName) (versions []*apitypes.ModrinthProjectVersion) {
	res, _ := http.Get(constructProjectVersionsUrl(slug))
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &versions)
	return
}

func GetProjectId(slug syntaxtypes.PackageName) (id string) {
	res, _ := http.Get(constructProjectUrl(slug))
	modrinthProject := apitypes.ModrinthProject{}
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &modrinthProject)
	id = modrinthProject.Id
	return
}

// TODO: Search() is way too long, refactor

func Search(
	platform syntaxtypes.Platform,
	packageName syntaxtypes.PackageName,
	showClientPackage bool,
	indexBy string,
) (result *apitypes.ModrinthSearchResults) {
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
	case syntaxtypes.AllPlatform:
		facetsArray = append(facetsArray, facetsCategoryAll)
	case syntaxtypes.Forge:
		facetsArray = append(facetsArray, facetsCategoryForge)
	case syntaxtypes.Fabric:
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

func modrinthProjectToPackageInfo(s *apitypes.ModrinthProject) *lucytypes.PackageInfo {
	name := syntaxtypes.PackageName(s.Slug)
	serverInfo := probe.GetServerInfo()

	info := &lucytypes.PackageInfo{
		Id: syntaxtypes.Package{
			// Edit in later code
			Platform: "",
			Name:     name,
			Version:  "",
		},
		Installed:          false,                    // Edit in later code
		Path:               "",                       // Edit in later code
		Urls:               []lucytypes.PackageUrl{}, // Edit in later code
		Name:               s.Title,
		Description:        s.Description,
		SupportedVersions:  []syntaxtypes.PackageVersion{}, // Edit in later code
		SupportedPlatforms: []syntaxtypes.Platform{},       // Edit in later code
	}

	// See if the mod is installed locally
	for _, mod := range serverInfo.Mods {
		if mod.Id.Name == name {
			info.Path = mod.Path
			info.Installed = true
			info.Id = mod.Id
		}
	}

	// Fill in supported versions and platforms
	for _, version := range s.GameVersions {
		info.SupportedVersions = append(
			info.SupportedVersions,
			syntaxtypes.PackageVersion(version),
		)
	}

	for _, platform := range s.Loaders {
		if syntaxtypes.Platform(platform).Valid() {
			info.SupportedPlatforms = append(
				info.SupportedPlatforms,
				syntaxtypes.Platform(platform),
			)
		}
	}

	// Fill in URLs
	if s.WikiUrl != "" {
		info.Urls = append(
			info.Urls, lucytypes.PackageUrl{
				Name: lucytypes.WikiUrl.String(),
				Type: lucytypes.WikiUrl,
				Url:  s.WikiUrl,
			},
		)
	}
	if s.SourceUrl != "" {
		info.Urls = append(
			info.Urls, lucytypes.PackageUrl{
				Name: lucytypes.HomepageUrl.String(),
				Type: lucytypes.SourceUrl,
				Url:  s.SourceUrl,
			},
		)
	}
	for _, url := range s.DonationUrls {
		info.Urls = append(
			info.Urls, lucytypes.PackageUrl{
				Name: "Donation",
				Type: lucytypes.OthersUrl,
				Url:  url.Url,
			},
		)
	}

	return info
}

package modrinth

import (
	"encoding/json"
	"html/template"
	"io"
	"lucy/apitypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/probe"
	"lucy/syntaxtypes"
	"net/http"
	"net/url"
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

const searchUrlTemplate = `https://api.modrinth.com/v2/search?query={{.packageName}}&limit=100&index={{.indexBy}}&facets={{.facets}}`

func Search(
	platform syntaxtypes.Platform,
	packageName syntaxtypes.PackageName,
	showClientPackage bool,
	indexBy string, // indexBy can be: relevance (default), downloads, follows, newest, updated
) (result *apitypes.ModrinthSearchResults) {
	// Construct the search url
	var facets []*facet
	switch platform {
	case syntaxtypes.Forge:
		facets = append(facets, facetForge)
	case syntaxtypes.Fabric:
		facets = append(facets, facetFabric)
	}

	if showClientPackage {
		facets = append(facets, facetBothSupported)
	} else {
		facets = append(facets, facetServerSupported)
	}

	templateUrl, _ := template.New("template_url").Parse(searchUrlTemplate)
	urlBuilder := strings.Builder{}
	_ = templateUrl.Execute(
		&urlBuilder,
		map[string]any{
			"packageName": string(packageName),
			"indexBy":     indexBy,
			"facets":      url.QueryEscape(StringifyFacets(facets...)),
		},
	)
	searchUrl := urlBuilder.String()

	// Make the call to Modrinth API
	println("calling", searchUrl)
	res, err := http.Get(searchUrl)
	if err != nil {
		logger.CreateFatal(err)
	}
	data, err := io.ReadAll(res.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.CreateWarning(err)
		}
	}(res.Body)
	err = json.Unmarshal(data, &result)

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
	for _, donationUrl := range s.DonationUrls {
		info.Urls = append(
			info.Urls, lucytypes.PackageUrl{
				Name: "Donation",
				Type: lucytypes.OthersUrl,
				Url:  donationUrl.Url,
			},
		)
	}

	return info
}

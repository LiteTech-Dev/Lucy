package modrinth

import (
	"encoding/json"
	"html/template"
	"io"
	"lucy/apitypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/probe"
	"net/http"
	"net/url"
	"strings"
)

// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/

func GetNewestProjectVersion(slug lucytypes.PackageName) (newestVersion *apitypes.ModrinthProjectVersion) {
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

func getProjectVersions(slug lucytypes.PackageName) (versions []*apitypes.ModrinthProjectVersion) {
	res, _ := http.Get(constructProjectVersionsUrl(slug))
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &versions)
	return
}

func GetProjectId(slug lucytypes.PackageName) (id string) {
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
	platform lucytypes.Platform,
	packageName lucytypes.PackageName,
	showClientPackage bool,
	indexBy string, // indexBy can be: relevance (default), downloads, follows, newest, updated
) (result *apitypes.ModrinthSearchResults) {
	// Construct the search url
	var facets []*facet
	switch platform {
	case lucytypes.Forge:
		facets = append(facets, facetForge)
	case lucytypes.Fabric:
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

func PackageFromModrinth(s *apitypes.ModrinthProject) *lucytypes.Package {
	p := &lucytypes.Package{}
	p.Dependencies = &lucytypes.PackageDependencies{}

	// Fill in supported versions and platforms
	for _, version := range s.GameVersions {
		p.Dependencies.SupportedVersions = append(
			p.Dependencies.SupportedVersions,
			lucytypes.PackageVersion(version),
		)
	}

	for _, platform := range s.Loaders {
		pf := lucytypes.Platform(platform)
		if pf.Valid() {
			p.Dependencies.SupportedPlatforms = append(
				p.Dependencies.SupportedPlatforms,
				pf,
			)
		}
	}

	p.Information = &lucytypes.PackageInformation{}

	// Fill in URLs
	if s.WikiUrl != "" {
		p.Information.Urls = append(
			p.Information.Urls, lucytypes.PackageUrl{
				Name: lucytypes.WikiUrl.String(),
				Type: lucytypes.WikiUrl,
				Url:  s.WikiUrl,
			},
		)
	}

	if s.SourceUrl != "" {
		p.Information.Urls = append(
			p.Information.Urls, lucytypes.PackageUrl{
				Name: lucytypes.HomepageUrl.String(),
				Type: lucytypes.SourceUrl,
				Url:  s.SourceUrl,
			},
		)
	}

	for _, donationUrl := range s.DonationUrls {
		p.Information.Urls = append(
			p.Information.Urls, lucytypes.PackageUrl{
				Name: "Donation",
				Type: lucytypes.OthersUrl,
				Url:  donationUrl.Url,
			},
		)
	}

	// Fill in the rest of the info
	p.Information.Brief = s.Description
	p.Information.Description = s.Body // s.Body is markdown, so it needs further processing
	p.Information.License = s.License.Name
	// p.Information.Author TODO: Author info

	return p
}

// Package modrinth provides functions to interact with Modrinth API
//
// We here use Modrinth terms in private functions:
//   - Project: A project is a mod, plugin, or resource pack.
//   - Version: A version is a release, beta, or alpha version of a project.
//
// Generally, a project in Modrinth is equivalent to a project in Lucy. And
// a version in Modrinth is equivalent to a package in Lucy.
//
// Here, while referring to a project in lucy, we would try to the term "slug"
// to refer to the project (or it's name).
package modrinth

import (
	"encoding/json"
	"errors"
	"io"
	"lucy/apitypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/tools"
	"net/http"
)

var ErrorInvalidAPIResponse = errors.New("invalid data from modrinth api")

func Fetch(id lucytypes.PackageId) (
	remote *lucytypes.PackageRemote,
	err error,
) {
	remote = &lucytypes.PackageRemote{
		Source:   lucytypes.Modrinth,
		RemoteId: getProjectId(id.Name),
	}

	fileUrl, filename, err := GetFile(id)
	if err == nil {
		remote.FileUrl = fileUrl
		remote.Filename = filename
	}

	return remote, err
}

// TODO: Incomplete

func Information(id lucytypes.PackageId) (
	information *lucytypes.PackageInformation,
	err error,
) {
	projcet := getProjectByName(id.Name)
	information = &lucytypes.PackageInformation{
		Name:        projcet.Title,
		Brief:       projcet.Description,
		Description: tools.MarkdownToPlainText(projcet.Body),
		Author:      nil,
		Urls:        []lucytypes.PackageUrl{},
		License:     projcet.License.Name,
	}

	// Fill in URLs
	if projcet.WikiUrl != "" {
		information.Urls = append(
			information.Urls,
			lucytypes.PackageUrl{
				Name: "Wiki",
				Type: lucytypes.WikiUrl,
				Url:  projcet.WikiUrl,
			},
		)
	}

	if projcet.SourceUrl != "" {
		information.Urls = append(
			information.Urls,
			lucytypes.PackageUrl{
				Name: "Source Code",
				Type: lucytypes.SourceUrl,
				Url:  projcet.SourceUrl,
			},
		)
	}

	if projcet.DonationUrls != nil {
		for _, donationUrl := range projcet.DonationUrls {
			information.Urls = append(
				information.Urls,
				lucytypes.PackageUrl{
					Name: "Donation",
					Type: lucytypes.OthersUrl,
					Url:  donationUrl.Url,
				},
			)
		}
	}

	return information, nil
}

// Dependencies from Modrinth API is extremely unreliable. A local check (if any
// files were downloaded) is recommended.
func Dependencies(packageId lucytypes.PackageId) (dependencies *lucytypes.PackageDependencies) {
	projcet := getProjectByName(packageId.Name)
	version, _ := getVersion(packageId)
	dependencies = &lucytypes.PackageDependencies{
		SupportedVersions:  []lucytypes.PackageVersion{},
		SupportedPlatforms: []lucytypes.Platform{},
		Required:           []lucytypes.PackageId{},
	}

	for _, version := range projcet.GameVersions {
		dependencies.SupportedVersions = append(
			dependencies.SupportedVersions,
			lucytypes.PackageVersion(version),
		)
	}

	for _, platform := range projcet.Loaders {
		dependencies.SupportedPlatforms = append(
			dependencies.SupportedPlatforms,
			lucytypes.Platform(platform),
		)
	}

	for _, dependency := range version.Dependencies {
		switch dependency.DependencyType {
		case apitypes.ModrinthVersionDependencyTypeIncompatible:
			d, err := DependencyToPackage(packageId, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Incompatible = append(
				dependencies.Incompatible,
				d,
			)
		case apitypes.ModrinthVersionDependencyTypeOptional:
			d, err := DependencyToPackage(packageId, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Optional = append(
				dependencies.Optional,
				d,
			)
		case apitypes.ModrinthVersionDependencyTypeRequired:
			d, err := DependencyToPackage(packageId, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Required = append(
				dependencies.Required,
				d,
			)
		}
	}

	return dependencies
}

// TODO:
// func Search(packageId lucytypes.PackageId, options SearchOptions) (result []*lucytypes.PackageId, err error)

// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/

func Search(
	packageId lucytypes.PackageId,
	showClientPackage bool,
) (result *apitypes.ModrinthSearchResults, err error) {
	var facets []facetItems
	query := packageId.Name

	switch packageId.Platform {
	case lucytypes.Forge:
		facets = append(facets, facetForge)
	case lucytypes.Fabric:
		facets = append(facets, facetFabric)
	}

	if showClientPackage {
		facets = append(facets, facetServerSupported, facetClientSupported)
	} else {
		facets = append(facets, facetServerSupported)
	}

	option := searchOptions{
		index:  byRelevance,
		facets: facets,
	}
	searchUrl := searchUrl(query, option)

	// Make the call to Modrinth API
	logger.Debug("searching via modrinth api: " + searchUrl)
	res, err := http.Get(searchUrl)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	defer tools.CloseReader(res.Body, logger.Warning)
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetProjectByName(packageName lucytypes.PackageName) (
	project *apitypes.ModrinthProject,
	err error,
) {
	res, err := http.Get(projectUrl(string(packageName)))
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

func FullPackage(s *apitypes.ModrinthProject) *lucytypes.Package {
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

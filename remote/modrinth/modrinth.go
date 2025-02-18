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
	"lucy/datatypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/tools"
	"net/http"
)

var ErrorInvalidAPIResponse = errors.New("invalid data from modrinth api")

// Search
//
// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/
func Search(
packageId lucytypes.PackageId,
options lucytypes.SearchOptions,
) (result *lucytypes.SearchResults, err error) {
	var facets []facetItems
	query := packageId.Name

	// switch packageId.Platform {
	// case lucytypes.Forge:
	// 	facets = append(facets, facetForge)
	// case lucytypes.Fabric:
	// 	facets = append(facets, facetFabric)
	// }

	if options.ShowClientPackage {
		facets = append(facets, facetServerSupported, facetClientSupported)
	} else {
		facets = append(facets, facetServerSupported)
	}

	internalOptions := searchOptions{
		index:  options.IndexBy.ToModrinth(),
		facets: facets,
	}
	searchUrl := searchUrl(
		query,
		internalOptions,
	)

	// Make the call to Modrinth API
	logger.Debug("searching via modrinth api: " + searchUrl)
	resp, err := http.Get(searchUrl)
	if err != nil {
		return nil, ErrorInvalidAPIResponse
	}
	data, err := io.ReadAll(resp.Body)
	defer tools.CloseReader(resp.Body, logger.Warning)
	var searchResults datatypes.ModrinthSearchResults
	err = json.Unmarshal(data, &searchResults)
	if err != nil {
		return nil, err
	}

	if searchResults.Hits == nil {
		return nil, nil
	}
	result = &lucytypes.SearchResults{}
	result.Results = make([]string, 0, len(searchResults.Hits))
	result.Source = lucytypes.Modrinth
	for _, hit := range searchResults.Hits {
		result.Results = append(result.Results, hit.Slug)
	}
	return result, nil
}

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

func Information(id lucytypes.PackageId) (
information *lucytypes.PackageInformation,
err error,
) {
	project := getProjectByName(id.Name)
	information = &lucytypes.PackageInformation{
		Name:        project.Title,
		Brief:       project.Description,
		Description: tools.MarkdownToPlainText(project.Body),
		Author:      nil,
		Urls:        []lucytypes.PackageUrl{},
		License:     project.License.Name,
	}

	// Fill in URLs
	if project.WikiUrl != "" {
		information.Urls = append(
			information.Urls,
			lucytypes.PackageUrl{
				Name: "Wiki",
				Type: lucytypes.WikiUrl,
				Url:  project.WikiUrl,
			},
		)
	}

	if project.SourceUrl != "" {
		information.Urls = append(
			information.Urls,
			lucytypes.PackageUrl{
				Name: "Source Code",
				Type: lucytypes.SourceUrl,
				Url:  project.SourceUrl,
			},
		)
	}

	if project.DonationUrls != nil {
		for _, donationUrl := range project.DonationUrls {
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

	// Fill in authors
	members := getProjectMembers(project.Id)
	for _, member := range members {
		information.Author = append(
			information.Author,
			lucytypes.PackageMember{
				Name:  member.User.Username,
				Role:  member.Role,
				Url:   userHomepageUrl(member.User.Id),
				Email: member.User.Email,
			},
		)
	}

	return information, nil
}

// Dependencies from Modrinth API is extremely unreliable. A local check (if any
// files were downloaded) is recommended.
func Dependencies(packageId lucytypes.PackageId) (dependencies *lucytypes.PackageDependencies) {
	project := getProjectByName(packageId.Name)
	version, _ := getVersion(packageId)
	dependencies = &lucytypes.PackageDependencies{
		SupportedVersions:  []lucytypes.PackageVersion{},
		SupportedPlatforms: []lucytypes.Platform{},
		Required:           []lucytypes.PackageId{},
	}

	for _, version := range project.GameVersions {
		dependencies.SupportedVersions = append(
			dependencies.SupportedVersions,
			lucytypes.PackageVersion(version),
		)
	}

	for _, platform := range project.Loaders {
		dependencies.SupportedPlatforms = append(
			dependencies.SupportedPlatforms,
			lucytypes.Platform(platform),
		)
	}

	for _, dependency := range version.Dependencies {
		switch dependency.DependencyType {
		case datatypes.ModrinthVersionDependencyTypeIncompatible:
			d, err := DependencyToPackage(packageId, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Incompatible = append(
				dependencies.Incompatible,
				d,
			)
		case datatypes.ModrinthVersionDependencyTypeOptional:
			d, err := DependencyToPackage(packageId, &dependency)
			if err != nil {
				logger.Warning(err)
				continue
			}
			dependencies.Optional = append(
				dependencies.Optional,
				d,
			)
		case datatypes.ModrinthVersionDependencyTypeRequired:
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

func GetProjectByName(packageName lucytypes.PackageName) (
project *datatypes.ModrinthProject,
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
	project = &datatypes.ModrinthProject{}
	err = json.Unmarshal(data, project)
	return
}

func FullPackage(s *datatypes.ModrinthProject) *lucytypes.Package {
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

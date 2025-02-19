package modrinth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"lucy/datatypes"
	"lucy/lucytypes"
)

func getProjectId(slug lucytypes.PackageName) (id string) {
	res, _ := http.Get(projectUrl(string(slug)))
	modrinthProject := datatypes.ModrinthProject{}
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &modrinthProject)
	id = modrinthProject.Id
	return
}

func getProjectById(id string) (project *datatypes.ModrinthProject) {
	res, _ := http.Get(projectUrl(id))
	data, _ := io.ReadAll(res.Body)
	project = &datatypes.ModrinthProject{}
	json.Unmarshal(data, project)
	return
}

func getProjectByName(slug lucytypes.PackageName) (project *datatypes.ModrinthProject) {
	res, _ := http.Get(projectUrl(string(slug)))
	data, _ := io.ReadAll(res.Body)
	project = &datatypes.ModrinthProject{}
	json.Unmarshal(data, project)
	return
}

func getProjectMembers(id string) (members []*datatypes.ModrinthMember) {
	res, _ := http.Get(projectMemberUrl(id))
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &members)
	return
}

var ErrorInvalidDependency = errors.New("invalid dependency")

func DependencyToPackage(
	depedent lucytypes.PackageId,
	dependency *datatypes.ModrinthVersionDependencies,
) (
	p lucytypes.PackageId,
	err error,
) {
	var version *datatypes.ModrinthVersion
	var project *datatypes.ModrinthProject

	// I don't see a case where a package would depend on a project on another
	// platform. So, we can safely assume that the platform of the dependent
	// package is the same as the platform of the dependency.
	p.Platform = depedent.Platform

	if dependency.VersionId != "" && dependency.ProjectId != "" {
		version = getVersionById(dependency.VersionId)
		project = getProjectById(dependency.ProjectId)
	} else if dependency.VersionId != "" {
		version = getVersionById(dependency.VersionId)
		project = getProjectById(version.ProjectId)
	} else if dependency.ProjectId != "" {
		project = getProjectById(dependency.ProjectId)
		// This is not safe, TODO: use better inference method
		version = latestVersion(lucytypes.PackageName(project.Slug))
		p.Version = lucytypes.LatestVersion
	} else {
		return p, ErrorInvalidDependency
	}

	p.Name = lucytypes.PackageName(project.Slug)
	p.Version = lucytypes.PackageVersion(version.VersionNumber)

	return p, nil
}

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/urfave/cli/v3"
	"io"
	"lucy/probe"
	"lucy/types"
	"lucy/util"
	"net/http"
	url2 "net/url"
	"reflect"
)

var SubcmdAdd = &cli.Command{
	Name:  "add",
	Usage: "Add new mods, plugins, or server modules",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Ignore version, dependency, and platform warnings",
			Value:   false,
		},
	},
	Action: ActionAdd,
}

// ActionAdd
// Now the strategy is:
//
//   - Most up to date
//   - Compatible with the server version
//   - Release version
//
// TODO: Version specification
func ActionAdd(_ context.Context, cmd *cli.Command) error {
	// TODO: Platform specification
	// TODO: Platform compatibility check

	platform, packageName := parsePackageSyntax(cmd.Args().First())
	serverInfo := probe.GetServerInfo()

	if !serverInfo.HasLucy {
		return errors.New("lucy is not installed, run `lucy init` before downloading mods")
	}

	if platform == "mcdr" && !serverInfo.HasMcdr {
		// TODO: Deal with this
		println("MCDR is not installed, cannot download MCDR plugins")
		return errors.New("no mcdr")
	} else if platform != "all" && platform != serverInfo.Executable.ModLoaderType {
		// TODO: Deal with this
		return errors.New("platform mismatch")
	}

	newestVersion := getNewestMorinthProjectVersion(packageName)
	file := util.DownloadFile(
		// Not sure how to deal with multiple files
		// As the motivation for publishers to provide multiple files is unclear
		// TODO: Maybe add a prompt to let the user choose
		newestVersion.Files[0].Url,
		"mod",
		newestVersion.Files[0].Filename,
	)

	util.InstallMod(file)

	return nil
}

func getNewestMorinthProjectVersion(slug string) (newestVersion types.ModrinthProjectVersion) {
	newestVersion = types.ModrinthProjectVersion{}
	versions := getModrinthProjectVersions(slug)
	serverInfo := probe.GetServerInfo()
	for _, version := range versions {
		for _, gameVersion := range version.GameVersions {
			// This if statement is a bit long and unreadable
			// TODO: Refactor
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

func constructModrinthProjectVersionsUrl(slug string) (url string) {
	url, _ = url2.JoinPath(
		"https://api.modrinth.com/v2/project",
		slug,
		"version",
	)
	return
}

func getModrinthProjectVersions(slug string) (versions []types.ModrinthProjectVersion) {
	res, _ := http.Get(constructModrinthProjectVersionsUrl(slug))
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &versions)
	return
}

func getModrinthProjectId(slug string) (id string) {
	res, _ := http.Get(constructModrinthProjectUrl(slug))
	modrinthProject := types.ModrinthProject{}
	data, _ := io.ReadAll(res.Body)
	json.Unmarshal(data, &modrinthProject)
	id = modrinthProject.Id
	return
}

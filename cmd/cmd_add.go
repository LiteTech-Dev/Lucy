package cmd

import (
	"context"
	"errors"
	"github.com/urfave/cli/v3"
	"lucy/modrinth"
	"lucy/probe"
	"lucy/syntax"
	"lucy/util"
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
	// TODO: Error handling

	_, p := syntax.Parse(cmd.Args().First())
	serverInfo := probe.GetServerInfo()

	if !serverInfo.HasLucy {
		return errors.New("lucy is not installed, run `lucy init` before downloading mods")
	}

	if p.Platform == syntax.Mcdr && !serverInfo.HasMcdr {
		// TODO: Deal with this
		println("MCDR is not installed, cannot download MCDR plugins")
		return errors.New("no mcdr")
	} else if p.Platform != syntax.AllPlatform && string(p.Platform) != serverInfo.Executable.ModLoaderType {
		// TODO: Deal with this
		return errors.New("platform mismatch")
	}

	newestVersion := modrinth.GetNewestProjectVersion(p.PackageName)
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

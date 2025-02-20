package cmd

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"
	"lucy/local"
	"lucy/logger"
	"lucy/lucyerrors"
	"lucy/lucytypes"
	"lucy/remote/modrinth"
	"lucy/syntax"
	"lucy/util"
)

var subcmdAdd = &cli.Command{
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
	Action: actionAdd,
}

// actionAdd
// The strategy is:
//
//   - Most up to date
//   - Compatible with the server version
//   - Release version
//
// TODO: Version specification
func actionAdd(_ context.Context, cmd *cli.Command) error {
	// TODO: Platform specification
	// TODO: Platform compatibility check
	// TODO: Error handling

	p := syntax.Parse(cmd.Args().First())
	serverInfo := local.GetServerInfo()

	if !serverInfo.HasLucy {
		return errors.New("lucy is not installed, run `lucy init` before downloading mods")
	}

	if serverInfo.Executable == local.UnknownExecutable {
		// Case where the server is not detected
		return errors.New("no executable found, `lucy add` requires a server in current directory")
	} else if p.Platform == lucytypes.Mcdr && serverInfo.Mcdr == nil {
		// Case where MCDR is not installed but the user wants to download MCDR plugins
		// TODO: Deal with this
		logger.Error(errors.New("no mcdr found, while mcdr plugins requested"))
		return nil
	} else if p.Platform != lucytypes.AllPlatform && p.Platform != serverInfo.Executable.Platform {
		// Case where the platform of the mod is different from the server
		// TODO: Deal with this
		logger.Error(errors.New("platform mismatch"))
		return nil
	}

	newestVersion := modrinth.LatestCompatibleVersion(p.Name)
	downloadFile, err := util.DownloadFile(
		// Not sure how to deal with multiple files
		// As the motivation for publishers to provide multiple files is unclear
		// TODO: Maybe add a prompt to let the user choose
		newestVersion.Files[0].Url,
		"mod",
		newestVersion.Files[0].Filename,
	)
	if err != nil {
		if errors.Is(err, lucyerrors.NoLucyError) {
			logger.Warning(err)
		} else {
			logger.Error(errors.New("failed at downloading: " + err.Error()))
			return nil
		}
	}

	util.MoveFile(downloadFile, serverInfo.ModPath)

	return nil
}

package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
	"lucy/logger"
)

// Frontend should change when user do not run the program in CLI
var Frontend = "cli"

// Each subcommand (and its action function) should be in its own file

// Cli is the main command for lucy
var Cli = &cli.Command{
	Name:   "lucy",
	Usage:  "The Minecraft server-side package manager",
	Action: addDefaultBehaviour(mainAction),
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "log-file",
			Aliases: []string{"l"},
			Usage:   "Output the path to logfile",
			Value:   false,
		},
	},
	Commands: []*cli.Command{
		subcmdStatus,
		subcmdInfo,
		subcmdSearch,
		subcmdAdd,
		subcmdInit,
	},
}

func mainAction(_ context.Context, cmd *cli.Command) error {
	// sources.MultiSourceDownload(
	// 	[]string{
	// 		"https://cdn.modrinth.com/data/P7dR8mSH/versions/nyAmoHlr/fabric-api-0.87.2%2B1.19.4.jar",
	// 		"https://mediafilez.forgecdn.net/files/4834/896/fabric-api-0.87.2%2B1.19.4.jar",
	// 	},
	// 	"fabric-api-0.87.2+1.19.4.jar",
	// )
	if cmd.Bool("log-file") {
		println("Log file at", logger.LogFile.Name())
	}
	return nil
}

// addDefaultBehaviour is a high-order function that takes a cli.ActionFunc and
// returns a cli.ActionFunc that prints help and exit when there's no args and flags.
//
// This function is not necessarily used for every action function, as some
// action functions are expected to have no args and flags. E.g., `lucy status`.
func addDefaultBehaviour(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 && len(cmd.FlagNames()) == 0 {
			cli.ShowAppHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

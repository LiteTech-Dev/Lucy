package cmd

import (
	"context"

	"lucy/tools"

	"github.com/urfave/cli/v3"
)

// Frontend should change when user do not run the program in CLI
// This is prepared for possible GUI implementation
var Frontend = "cli"

// Each subcommand (and its action function) should be in its own file

// Cli is the main command for lucy
var Cli = &cli.Command{
	Name:  "lucy",
	Usage: "The Minecraft server-side package manager",
	Action: tools.Decorate(
		actionEmpty,
		globalFlagsDecorator,
		helpOnNoInputDecorator,
	),
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "log-file",
			Aliases: []string{"l"},
			Usage:   "Output the path to logfile",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Print logs",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Print debug logs",
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

var actionEmpty cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {
	return nil
}

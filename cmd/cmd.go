package cmd

import (
	"context"

	"lucy/tools"

	"github.com/urfave/cli/v3"
	"lucy/logger"
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

// globalFlagsDecorator is a high-order function that appends global flag actions
// to the action function.
func globalFlagsDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 && len(cmd.FlagNames()) == 0 {
			cli.ShowAppHelpAndExit(cmd, 0)
		}
		if cmd.Bool("log-file") {
			println("Log file at", logger.LogFile.Name())
		}
		if cmd.Bool("verbose") {
			logger.UseConsoleOutput()
		}
		if cmd.Bool("debug") {
			logger.UseDebug()
		}
		return f(ctx, cmd)
	}
}

// helpOnNoInputDecorator is a high-order function that takes a cli.ActionFunc and
// returns a cli.ActionFunc that prints help and exit when there's no args and flags.
//
// This function is not necessarily used for every action function, as some
// action functions are expected to have no args and flags. E.g., `lucy status`.
func helpOnNoInputDecorator(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 && len(cmd.LocalFlagNames()) == 0 {
			cli.ShowAppHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"lucy/output"
)

// Frontend
// This changes when user runs the web interface
var Frontend = "cli"

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
	if cmd.Bool("log-file") {
		println("Log file at", output.LogFile.Name())
	}
	return nil
}

// addDefaultBehaviour is a high-order function that takes an action func and returns
// an action func that prints help and exit when there's no args and flags.
func addDefaultBehaviour(f cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() == 0 && len(cmd.FlagNames()) == 0 {
			cli.ShowAppHelpAndExit(cmd, 0)
		}
		return f(ctx, cmd)
	}
}

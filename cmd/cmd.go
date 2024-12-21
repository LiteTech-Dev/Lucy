package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
)

// Frontend
// This changes when user runs the web interface
var Frontend = "cli"

var Cli = &cli.Command{
	Name:   "lucy",
	Usage:  "The Minecraft server-side package manager",
	Action: noArgAction,
	Commands: []*cli.Command{
		SubcmdStatus,
		SubcmdInfo,
		SubcmdSearch,
		SubcmdAdd,
		SubcmdInit,
	},
}

// This shows the help message of the called command
func noArgAction(_ context.Context, cmd *cli.Command) error {
	cli.ShowAppHelpAndExit(cmd, 0)
	return nil
}

package cli

import (
	"context"
	"github.com/urfave/cli/v3"
)

var Cli = &cli.Command{
	Commands: []*cli.Command{
		{},
	},
	Name:  "lucy",
	Usage: "The Minecraft server-side package manager",
	Action: func(ctx context.Context, command *cli.Command) error {
		cli.ShowAppHelpAndExit(command, 0)
		return nil
	},
}

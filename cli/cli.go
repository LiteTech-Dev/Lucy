package cli

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/probe"
)

var Cli = &cli.Command{
	Name:   "lucy",
	Usage:  "The Minecraft server-side package manager",
	Action: noArgAction,
	Commands: []*cli.Command{
		{
			Name:    "add",
			Usage:   "Add new mods, plugins, or server modules",
			Aliases: []string{"a"},
			Action:  noArgAction,
			Commands: []*cli.Command{
				{
					Name:   "mod",
					Usage:  "Add new mod(s)",
					Action: addMod,
				},
			},
		},
		{
			Name:    "list",
			Usage:   "List everything installed on your server",
			Aliases: []string{"ls"},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println(probe.GetServerFiles().ServerWorkPath)
				return nil
			},
		},
	},
}

// This shows the help message of the called command
func noArgAction(_ context.Context, cmd *cli.Command) error {
	cli.ShowAppHelpAndExit(cmd, 0)
	return nil
}

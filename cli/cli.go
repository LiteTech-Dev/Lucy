package cli

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/probe"
	"strings"
)

var Frontend = "cli"
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
					Name:  "mod",
					Usage: "Add new mod(s)",
				},
			},
		},
		{
			Name:    "info",
			Usage:   "Display information of the current server",
			Aliases: []string{"i"},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				serverInfo := probe.GetServerInfo()
				fmt.Printf("Minecraft v%s\n", serverInfo.Executable.GameVersion)
				fmt.Printf(
					"%s%s ",
					strings.ToUpper(serverInfo.Executable.ModLoaderType[:1]),
					serverInfo.Executable.ModLoaderType[1:],
				)
				if serverInfo.Executable.ModLoaderType != "vanilla" {
					fmt.Printf("v%s\n", serverInfo.Executable.ModLoaderVersion)
				} else {
					fmt.Printf("\n")
				}
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

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
			Name:  "search",
			Usage: "Search for mods and plugins",
			Flags: []cli.Flag{
				// TODO: This flag is not yet implemented
				&cli.StringFlag{
					Name:     "source",
					Aliases:  []string{"s"},
					Usage:    "To search from `SOURCE`",
					Value:    "modrinth",
					Required: false,
					Validator: func(s string) error {
						if s != "modrinth" && s != "curseforge" {
							return fmt.Errorf("unsupported source: %s", s)
						}
						return nil
					},
				},
			},
			Action: SubcmdSearch,
		},
		{
			Name:    "add",
			Usage:   "Add new mods, plugins, or server modules",
			Aliases: []string{"a"},
			Action:  noArgAction,
		},
		{
			Name:  "status",
			Usage: "Display basic information of the current server",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				serverInfo := probe.GetServerInfo()

				// Print game version
				fmt.Printf("Minecraft v%s\n", serverInfo.Executable.GameVersion)

				// Print MCDR status
				if serverInfo.HasMcdr {
					fmt.Printf("Managed by MCDR\n")
				}

				// Print mod loader type and version
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

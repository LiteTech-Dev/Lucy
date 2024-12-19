package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
)

var SubcmdAdd = &cli.Command{
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
	Action: noArgAction,
}

func ActionAdd(ctx context.Context, cmd *cli.Command) error {
	return nil
}

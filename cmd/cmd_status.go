package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"lucy/output"
	"lucy/probe"
	"lucy/tools"
)

var subcmdStatus = &cli.Command{
	Name:   "status",
	Usage:  "Display basic information of the current server",
	Action: actionStatus,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Output in raw JSON format",
		},
	},
}

func actionStatus(_ context.Context, cmd *cli.Command) error {
	serverInfo := probe.GetServerInfo()
	if cmd.Bool("debug") {
		tools.PrintAsJson(serverInfo)
	} else {
		output.GenerateStatus(&serverInfo)
	}
	return nil
}

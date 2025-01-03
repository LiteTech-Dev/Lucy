package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/output"
	"lucy/probe"
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
		data, _ := json.MarshalIndent(serverInfo, "", "  ")
		fmt.Println(string(data))
	} else {
		output.GenerateStatus(&serverInfo)
	}
	return nil
}

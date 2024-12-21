package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/output"
	"lucy/probe"
)

var SubcmdStatus = &cli.Command{
	Name:   "status",
	Usage:  "Display basic information of the current server",
	Action: ActionStatus,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Output in raw JSON format",
		},
	},
}

func ActionStatus(_ context.Context, cmd *cli.Command) error {
	serverInfo := probe.GetServerInfo()
	if cmd.Bool("debug") {
		data, _ := json.MarshalIndent(serverInfo, "", "  ")
		fmt.Println(string(data))
	} else {
		output.GenerateStatus(serverInfo)
	}
	return nil
}

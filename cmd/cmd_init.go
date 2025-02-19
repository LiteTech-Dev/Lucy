package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

var subcmdInit = &cli.Command{
	Name:   "init",
	Usage:  "Initialize Lucy on current directory",
	Action: ActionInit,
}

func ActionInit(ctx context.Context, cmd *cli.Command) error {
	return nil
}

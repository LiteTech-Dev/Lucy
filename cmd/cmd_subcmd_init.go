package cmd

import (
	"context"
	"github.com/urfave/cli/v3"
	"lucy/util"
	"os"
)

var SubcmdInit = &cli.Command{
	Name:   "init",
	Usage:  "Initialize Lucy on current directory",
	Action: InitAction,
}

func InitAction(ctx context.Context, cmd *cli.Command) error {
	os.Mkdir(util.LucyPath, 0755)
	return nil
}

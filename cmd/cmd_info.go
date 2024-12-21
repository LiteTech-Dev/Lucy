package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"lucy/mcdr"
	"lucy/modrinth"
	"lucy/output"
	"lucy/syntax"
	"lucy/types"
	"net/http"
)

var SubcmdInfo = &cli.Command{
	Name:  "info",
	Usage: "Display information of a mod or plugin",
	Flags: []cli.Flag{
		// TODO: This flag is not yet implemented
		&cli.StringFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Usage:   "To fetch info from `SOURCE`",
			Value:   "modrinth",
		},
		&cli.BoolFlag{
			Name:    "markdown",
			Aliases: []string{"Md"},
			Usage:   "Print raw Markdown",
			Value:   false,
		},
	},
	Action: ActionInfo,
}

func ActionInfo(ctx context.Context, cmd *cli.Command) error {
	// TODO: Error handling
	_, p := syntax.Parse(cmd.Args().First())

	switch p.Platform {
	case syntax.AllPlatform:
		// TODO: Wide range search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.PackageName))
		modrinthProject := &types.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, modrinthProject)
		output.GenerateInfo(modrinthProject)
	case syntax.Fabric:
		// TODO: Fabric specific search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.PackageName))
		modrinthProject := &types.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, modrinthProject)
		output.GenerateInfo(modrinthProject)
	case syntax.Forge:
		// TODO: Forge support
		println("Not yet implemented")
	case syntax.Mcdr:
		mcdrPlugin := mcdr.SearchMcdrPluginCatalogue(p.PackageName)
		if mcdrPlugin == nil {
			_ = fmt.Errorf("plugin not found")
			return nil
		}
		output.GenerateInfo(mcdrPlugin)
	}

	return nil
}

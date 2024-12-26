package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"lucy/apitypes"
	"lucy/mcdr"
	"lucy/output"
	"lucy/sources/modrinth"
	"lucy/syntax"
	"net/http"
)

var subcmdInfo = &cli.Command{
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
	Action: actionInfo,
}

func actionInfo(ctx context.Context, cmd *cli.Command) error {
	// TODO: Error handling
	p, _ := syntax.Parse(cmd.Args().First())

	switch p.Platform {
	case syntax.AllPlatform:
		// TODO: Wide range search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.Name))
		modrinthProject := &apitypes.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, modrinthProject)
		output.GenerateInfo(modrinthProject)
	case syntax.Fabric:
		// TODO: Fabric specific search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.Name))
		modrinthProject := &apitypes.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, modrinthProject)
		output.GenerateInfo(modrinthProject)
	case syntax.Forge:
		// TODO: Forge support
		println("Not yet implemented")
	case syntax.Mcdr:
		mcdrPlugin := mcdr.SearchMcdrPluginCatalogue(p.Name)
		if mcdrPlugin == nil {
			_ = fmt.Errorf("plugin not found")
			return nil
		}
		output.GenerateInfo(mcdrPlugin)
	}

	return nil
}

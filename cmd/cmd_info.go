package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"lucy/apitypes"
	"lucy/mcdr"
	"lucy/modrinth"
	"lucy/output"
	"lucy/syntax"
	"lucy/syntaxtypes"
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
	p := syntax.Parse(cmd.Args().First())

	switch p.Platform {
	case syntaxtypes.AllPlatform:
		// TODO: Wide range search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.Name))
		modrinthProject := &apitypes.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		err := json.Unmarshal(data, modrinthProject)
		if err != nil {
			return err
		}
		output.GenerateInfo(modrinthProject)
	case syntaxtypes.Fabric:
		// TODO: Fabric specific search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.Name))
		modrinthProject := &apitypes.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		err := json.Unmarshal(data, modrinthProject)
		if err != nil {
			return err
		}
		output.GenerateInfo(modrinthProject)
	case syntaxtypes.Forge:
		// TODO: Forge support
		println("Not yet implemented")
	case syntaxtypes.Mcdr:
		mcdrPlugin := mcdr.SearchMcdrPluginCatalogue(p.Name)
		if mcdrPlugin == nil {
			_ = fmt.Errorf("plugin not found")
			return nil
		}
		output.GenerateInfo(mcdrPlugin)
	}

	return nil
}

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"lucy/apitypes"
	"lucy/lucytypes"
	"lucy/mcdr"
	"lucy/modrinth"
	"lucy/output"
	"lucy/syntax"
	"lucy/syntaxtypes"
	"lucy/tools"
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

	var multiSourceData []*lucytypes.OutputData

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
		multiSourceData = append(
			multiSourceData,
			modrinthProjectToInfo(modrinthProject),
		)
	case syntaxtypes.Fabric:
		// TODO: Fabric specific search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.Name))
		modrinthProject := &apitypes.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		err := json.Unmarshal(data, modrinthProject)
		if err != nil {
			return err
		}
		multiSourceData = append(
			multiSourceData,
			modrinthProjectToInfo(modrinthProject),
		)
	case syntaxtypes.Forge:
		// TODO: Forge support
		println("Not yet implemented")
	case syntaxtypes.Mcdr:
		mcdrPlugin := mcdr.SearchMcdrPluginCatalogue(p.Name)
		if mcdrPlugin == nil {
			_ = fmt.Errorf("plugin not found")
			return nil
		}
		multiSourceData = append(
			multiSourceData,
			mcdrPluginInfoToInfo(mcdrPlugin),
		)
	}

	for _, data := range multiSourceData {
		output.GenerateOutput(data)
	}

	return nil
}

func modrinthProjectToInfo(source *apitypes.ModrinthProject) *lucytypes.OutputData {
	return &lucytypes.OutputData{
		Fields: []lucytypes.Field{
			&output.FieldShortText{
				Title: "Name",
				Text:  source.Title,
			},
			&output.FieldShortText{
				Title: "Description",
				Text:  source.Description,
			},
			&output.FieldShortText{
				Title: "Downloads",
				Text:  fmt.Sprintf("%d", source.Downloads),
			},
			&output.FieldLabels{
				Title:    "Versions",
				Labels:   source.GameVersions,
				MaxWidth: 0,
			},
		},
	}
}

func mcdrPluginInfoToInfo(source *apitypes.McdrPluginInfo) *lucytypes.OutputData {
	return &lucytypes.OutputData{
		Fields: []lucytypes.Field{
			&output.FieldShortText{
				Title: "Name",
				Text:  source.Id,
			},
			&output.FieldShortText{
				Title: "Introduction",
				Text:  source.Introduction.EnUs,
			},
			&output.FieldPeople{
				Title: "Authors",
				People: []struct {
					Name string
					Link string
				}(source.Authors),
			},
			&output.FieldShortText{
				Title: "Source Code",
				Text:  tools.Underline(source.Repository),
			},
		},
	}
}

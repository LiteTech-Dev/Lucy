package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/apitypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/mcdr"
	"lucy/output"
	"lucy/sources/modrinth"
	"lucy/syntax"
	"lucy/syntaxtypes"
	"lucy/tools"
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
		proj, err := modrinth.GetProjectByName(p.Name)
		if err != nil {
			logger.CreateWarning(err)
			break
		}
		multiSourceData = append(
			multiSourceData,
			modrinthProjectToInfo(proj),
		)
	case syntaxtypes.Fabric:
		// TODO: Fabric specific search
		proj, err := modrinth.GetProjectByName(p.Name)
		if err != nil {
			logger.CreateWarning(err)
			break
		}
		multiSourceData = append(
			multiSourceData,
			modrinthProjectToInfo(proj),
		)
	case syntaxtypes.Forge:
		// TODO: Forge
		logger.CreateFatal(fmt.Errorf("forge is not yet supported"))
	case syntaxtypes.Mcdr:
		mcdrPlugin, err := mcdr.SearchMcdrPluginCatalogue(p.Name)
		if err != nil {
			logger.CreateWarning(err)
			break
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

// TODO: Link to newest version
// TODO: Link to latest compatible version
// TODO: Generate `lucy install` command

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

// TODO: Link to newest version
// TODO: Generate `lucy install` command

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

package cmd

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"lucy/datatypes"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/remote/mcdr"
	"lucy/remote/modrinth"
	"lucy/syntax"
	"lucy/tools"

	"github.com/urfave/cli/v3"
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
	Action: tools.Decorate(
		actionInfo,
		globalFlagsDecorator,
		helpOnNoInputDecorator,
	),
}

var actionInfo cli.ActionFunc = func(
	ctx context.Context,
	cmd *cli.Command,
) error {
	p := syntax.Parse(cmd.Args().First())

	var multiSourceData []*lucytypes.OutputData

	switch p.Platform {
	case lucytypes.AllPlatform:
		var packageFromModrinth lucytypes.Package
		packageFromModrinth.Remote, _ = modrinth.Fetch(p)
		packageFromModrinth.Information, _ = modrinth.Information(p.Name)
		packageFromModrinth.Dependencies = modrinth.Dependencies(p)
		multiSourceData = append(
			multiSourceData,
			cInfoOutput(packageFromModrinth),
		)
	case lucytypes.Fabric:
		// TODO: Fabric specific search
		modrinthProject, err := modrinth.GetProjectByName(p.Name)
		if err != nil {
			logger.Warning(err)
			break
		}
		multiSourceData = append(
			multiSourceData,
			modrinthProjectToInfo(modrinthProject),
		)
	case lucytypes.Forge:
		// TODO: Forge
		logger.Fatal(fmt.Errorf("forge is not yet supported"))
	case lucytypes.Mcdr:
		mcdrPlugin, err := mcdr.SearchMcdrPluginCatalogue(p.Name)
		if err != nil {
			logger.Warning(err)
			break
		}
		multiSourceData = append(
			multiSourceData,
			mcdrPluginInfoToInfo(mcdrPlugin),
		)
	}

	for _, data := range multiSourceData {
		output.Flush(data)
	}

	return nil
}

// TODO: Link to newest version
// TODO: Link to latest compatible version
// TODO: Generate `lucy add` command

func modrinthProjectToInfo(source *datatypes.ModrinthProject) *lucytypes.OutputData {
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
				Text:  strconv.Itoa(source.Downloads),
			},
			&output.FieldLabels{
				Title:    "Versions",
				Labels:   source.GameVersions,
				MaxWidth: 0,
			},
		},
	}
}

func mcdrPluginInfoToInfo(source *datatypes.McdrPluginInfo) *lucytypes.OutputData {
	info := &lucytypes.OutputData{
		Fields: []lucytypes.Field{
			&output.FieldShortText{
				Title: "Name",
				Text:  source.Id,
			},
			&output.FieldShortText{
				Title: "Description",
				Text:  source.Introduction.EnUs,
			},
			&output.FieldMultiShortTextWithAnnot{
				Title:  "Authors",
				Texts:  []string{},
				Annots: []string{},
			},
			&output.FieldShortText{
				Title: "Source Code",
				Text:  tools.Underline(source.Repository),
			},
		},
	}

	// This is temporary TODO: Use iota for fields instead
	const authorsField = 2
	a := info.Fields[authorsField].(*output.FieldMultiShortTextWithAnnot)

	for _, p := range source.Authors {
		a.Texts = append(a.Texts, p.Name)
		a.Annots = append(a.Annots, tools.Underline(p.Link))
	}

	return info
}

func cInfoOutput(p lucytypes.Package) *lucytypes.OutputData {
	o := &lucytypes.OutputData{
		Fields: []lucytypes.Field{
			&output.FieldAnnotation{
				Annotation: "(from " + p.Remote.Source.String() + ")",
			},
			&output.FieldShortText{
				Title: "Name",
				Text:  p.Information.Name,
			},
			&output.FieldShortText{
				Title: "Description",
				Text:  p.Information.Brief,
			},
		},
	}

	var authorNames []string
	var authorLinks []string
	for _, author := range p.Information.Author {
		authorNames = append(authorNames, author.Name)
		// TODO: Improve author info annotation format
		authorLinks = append(authorLinks, author.Url)
	}

	for _, url := range p.Information.Urls {
		o.Fields = append(
			o.Fields, &output.FieldShortText{
				Title: url.Name,
				Text:  tools.Underline(url.Url),
			},
		)
	}

	o.Fields = append(
		o.Fields, &output.FieldAnnotatedShortText{
			Title:      "Download",
			Text:       tools.Underline(p.Remote.FileUrl),
			Annotation: p.Remote.Filename,
		},
	)

	// TODO: Put current server version on the top
	// TODO: Hide snapshot versions, except if the current server is using it
	if !slices.Contains(p.Dependencies.SupportedPlatforms, lucytypes.Mcdr) &&
		(p.Dependencies.SupportedPlatforms != nil || len(p.Dependencies.SupportedPlatforms) != 0) {
		f := &output.FieldLabels{
			Title:    "Game Versions",
			Labels:   []string{},
			MaxWidth: 0,
			MaxLines: tools.TermHeight() / 2,
		}
		for _, version := range p.Dependencies.SupportedVersions {
			f.Labels = append(f.Labels, version.String())
		}
		o.Fields = append(o.Fields, f)
	}

	return o
}

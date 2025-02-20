package cmd

import (
	"context"
	"errors"
	"strconv"

	"github.com/urfave/cli/v3"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/remote/modrinth"
	"lucy/syntax"
	"lucy/tools"
)

var subcmdSearch = &cli.Command{
	Name:  "search",
	Usage: "Search for mods and plugins",
	Flags: []cli.Flag{
		// TODO: This flag is not yet implemented
		sourceFlag(lucytypes.Modrinth),
		&cli.StringFlag{
			Name:    "index",
			Aliases: []string{"i"},
			Usage:   "Index search results by `INDEX`",
			Value:   "relevance",
			Validator: func(s string) error {
				if lucytypes.SearchIndex(s).Validate() {
					return nil
				}
				return errors.New("must be one of \"relevance\", \"downloads\",\"newest\"")
			},
		},
		&cli.BoolFlag{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "Also show client-only mods in results",
			Value:   false,
		},
		flagJsonOutput,
		flagLongOutput,
	},
	Action: tools.Decorate(
		actionSearch,
		globalFlagsDecorator,
		helpOnNoInputDecorator,
	),
}

var actionSearch cli.ActionFunc = func(
_ context.Context,
cmd *cli.Command,
) error {
	p := syntax.Parse(cmd.Args().First())
	_ = cmd.String("index")
	showClientPackage := cmd.Bool("client")
	indexBy := lucytypes.SearchIndex(cmd.String("index"))

	res, err := modrinth.Search(
		p,
		lucytypes.SearchOptions{
			ShowClientPackage: showClientPackage,
			IndexBy:           indexBy,
		},
	)
	if err != nil {
		logger.Fatal(err)
	}
	output.Flush(generateSearchOutput(res, cmd.Bool("long")))

	return nil
}

func generateSearchOutput(
res *lucytypes.SearchResults,
showAll bool,
) *lucytypes.OutputData {
	return &lucytypes.OutputData{
		Fields: []lucytypes.Field{
			&output.FieldShortText{
				Title: "#  ",
				Text:  strconv.Itoa(len(res.Results)),
			},
			&output.FieldDynamicColumnLabels{
				Title:    ">>>",
				Labels:   res.Results,
				MaxLines: tools.Ternary(showAll, 0, tools.TermHeight()-6),
			},
		},
	}
}

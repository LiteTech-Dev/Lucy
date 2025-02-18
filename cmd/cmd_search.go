package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/logger"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/remote/modrinth"
	"lucy/syntax"
	"lucy/tools"
	"strconv"
)

var subcmdSearch = &cli.Command{
	Name:  "search",
	Usage: "Search for mods and plugins",
	Flags: []cli.Flag{
		// TODO: This flag is not yet implemented
		&cli.StringFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Usage:   "To search from `SOURCE`",
			Value:   "modrinth",
			Validator: func(s string) error {
				if s != "modrinth" && s != "curseforge" {
					return fmt.Errorf("unsupported source: %s", s)
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:    "index",
			Aliases: []string{"i"},
			Usage:   "Index search results by `INDEX`",
			Value:   "relevance",
			Validator: func(s string) error {
				if s != "relevance" && s != "downloads" && s != "follows" && s != "newest" && s != "updated" {
					return fmt.Errorf(
						`unsupported index: %s, value must be one of "relevance", "downloads", "follows", "newest", "updated"`,
						s,
					)
				}
				return nil
			},
		},
		&cli.BoolFlag{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "Also show client-only mods in results",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Output raw JSON response",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "all",
			Usage:   "Show all search results",
			Value:   false,
			Aliases: []string{"a"},
		},
	},
	Action: actionSearch,
}

func actionSearch(_ context.Context, cmd *cli.Command) error {
	p := syntax.Parse(cmd.Args().First())
	_ = cmd.String("index")
	showClientPackage := cmd.Bool("client")
	indexBy := lucytypes.InputSearchIndex(cmd.String("index"))

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
	output.Flush(generateSearchOutput(res, cmd.Bool("all")))

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
				MaxLines: tools.Ternary(showAll, 0, tools.TermHeight()-6
				),
			},
		},
	}
}

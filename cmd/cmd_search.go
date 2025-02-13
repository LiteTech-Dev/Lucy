package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/apitypes"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/sources/modrinth"
	"lucy/syntax"
	"lucy/tools"
	"os"
	"strconv"
	"text/tabwriter"
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
	},
	Action: actionSearch,
}

func actionSearch(_ context.Context, cmd *cli.Command) error {
	// TODO: Error handling
	p := syntax.Parse(cmd.Args().First())
	// indexBy can be: relevance (default), downloads, follows, newest, updated
	indexBy := cmd.String("index")
	showClientPackage := cmd.Bool("client")

	res := modrinth.Search(
		p.Platform,
		p.Name,
		showClientPackage,
		indexBy,
	)

	if cmd.Bool("debug") {
		tools.PrintAsJson(res)
		return nil
	}
	output.Flush(modrinthResToSearch(res))
	return nil
}

func modrinthResToSearch(res *apitypes.ModrinthSearchResults) *lucytypes.OutputData {
	hits := make([]string, len(res.Hits))
	for i, hit := range res.Hits {
		hits[i] = hit.Slug
	}
	return &lucytypes.OutputData{
		Fields: []lucytypes.Field{
			&output.FieldShortText{
				Title: "#  ",
				Text:  strconv.Itoa(res.TotalHits),
			},
			&output.FieldDynamicColumnLabels{
				Title:  ">>>",
				Labels: hits,
			},
		},
	}
}

func generateSearchOutput(slugs []lucytypes.PackageName) {
	maxSlugLen := 0
	for i := 0; i < len(slugs); i += 1 {
		if len(slugs[i]) > maxSlugLen {
			maxSlugLen = len(slugs[i])
		}
	}
	columns := tools.TermWidth() / (maxSlugLen + 2)

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Printf("Found %d results from Modrinth\n", len(slugs))
	for i := 0; i < len(slugs); i += 1 {
		if (i+1)%columns == 0 || i == len(slugs)-1 {
			fmt.Fprintf(writer, "%s\n", slugs[i])
		} else {
			fmt.Fprintf(writer, "%s\t", slugs[i])
		}
	}

	writer.Flush()
}

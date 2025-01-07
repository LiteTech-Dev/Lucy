package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
	"lucy/modrinth"
	"lucy/syntax"
	"os"
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
	_, p := syntax.Parse(cmd.Args().First())
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
		jsonOutput, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(jsonOutput))
		return nil
	}

	var slugs []syntax.PackageName
	for _, hit := range res.Hits {
		slugs = append(slugs, syntax.PackageName(hit.Slug))
	}
	generateSearchOutput(slugs)

	return nil
}

func generateSearchOutput(slugs []syntax.PackageName) {
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	maxSlugLen := 0
	for i := 0; i < len(slugs); i += 1 {
		if len(slugs[i]) > maxSlugLen {
			maxSlugLen = len(slugs[i])
		}
	}
	columns := termWidth / (maxSlugLen + 2)

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

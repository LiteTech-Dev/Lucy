package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
	"html/template"
	"io"
	"lucy/types"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

var SubcmdSearch = &cli.Command{
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
	Action: ActionSearch,
}

func ActionSearch(_ context.Context, cmd *cli.Command) error {
	platform, packageName := parsePackageSyntax(cmd.Args().First())
	// indexBy can be: relevance (default), downloads, follows, newest, updated
	indexBy := cmd.String("index")
	showClientPackage := cmd.Bool("client")

	res := searchModrinth(platform, packageName, showClientPackage, indexBy)
	if cmd.Bool("debug") {
		jsonOutput, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(jsonOutput))
		return nil
	}

	var slugs []string
	for _, hit := range res.Hits {
		slugs = append(slugs, hit.Slug)
	}
	generateSearchOutput(slugs)

	return nil
}

// For Modrinth search API, see:
// https://docs.modrinth.com/api/operations/searchprojects/

func searchModrinth(
	platform string,
	packageName string,
	showClientPackage bool,
	indexBy string,
) (result types.ModrinthSearchRes) {
	// Construct the search url
	const facetsCategoryAll = `["categories:'forge'","categories:'fabric'","categories:'quilt'","categories:'liteloader'","categories:'modloader'","categories:'rift'","categories:'neoforge'"]`
	const facetsCategoryForge = `["categories:'forge'"]`
	const facetsCategoryFabric = `["categories:'fabric'"]`
	const facetsServerOnly = `["server_side:optional","server_side:required"],["client_side:optional","client_side:required","client_side:unsupported"]`
	const facetsShowClient = `["server_side:optional","server_side:required","server_side:unsupported"],["client_side:optional","client_side:required","client_side:unsupported"]`
	const facetsTypeMod = `["project_type:'mod'"]`
	const urlTemplate = `https://api.modrinth.com/v2/search?query={{.packageName}}&limit=100&index={{.indexBy}}&facets={{.facetsEncoded}}`

	var facetsArray []string
	switch platform {
	case "all":
		facetsArray = append(facetsArray, facetsCategoryAll)
	case "forge":
		facetsArray = append(facetsArray, facetsCategoryForge)
	case "fabric":
		facetsArray = append(facetsArray, facetsCategoryFabric)
	}
	if !showClientPackage {
		facetsArray = append(facetsArray, facetsServerOnly)
	} else {
		facetsArray = append(facetsArray, facetsShowClient)
	}
	facetsArray = append(facetsArray, facetsTypeMod)
	facetsEncoded := url.QueryEscape("[" + strings.Join(facetsArray, ",") + "]")

	templateUrl, _ := template.New("template_url").Parse(urlTemplate)
	urlBuilder := strings.Builder{}
	_ = templateUrl.Execute(
		&urlBuilder,
		map[string]string{
			"packageName":   packageName,
			"indexBy":       indexBy,
			"facetsEncoded": facetsEncoded,
		},
	)

	// Make the call to Modrinth API
	resp, _ := http.Get(urlBuilder.String())
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	resp.Body.Close()

	return
}

func generateSearchOutput(slugs []string) {
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	maxSlugLen := 0
	for i := 0; i < len(slugs); i += 1 {
		if len(slugs[i]) > maxSlugLen {
			maxSlugLen = len(slugs[i])
		}
	}
	columns := termWidth / (maxSlugLen + 2)

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Printf("Found %d results\n", len(slugs))
	for i := 0; i < len(slugs); i += 1 {
		if (i+1)%columns == 0 || i == len(slugs)-1 {
			fmt.Fprintf(writer, "%s\n", slugs[i])
		} else {
			fmt.Fprintf(writer, "%s\t", slugs[i])
		}
	}

	writer.Flush()
}

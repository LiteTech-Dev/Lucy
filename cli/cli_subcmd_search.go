package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
	"html/template"
	"io"
	"lucy/types"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

// SubcmdSearch
// Search syntax:
// lucy search <query>
// Query can either be a single word or a string in the format of "loader/name"
// Example: lucy search carpet
// Example: lucy search fabric/carpet
// Example: lucy search mcdr/prime-backup
func SubcmdSearch(ctx context.Context, cmd *cli.Command) error {
	query := strings.Split(cmd.Args().First(), "/")
	res := types.ModrinthSearchRes{}
	if len(query) == 1 {
		res = searchModrinth(query[0], "fabric")
	} else if len(query) == 2 {
		// TODO: Implement categorized syntax
		println("Not yet implemented")
	} else {
		return errors.New("invalid query format")
	}

	outputs := []string{}
	for _, hit := range res.Hits {
		outputs = append(outputs, hit.Slug)
	}
	fmt.Printf("Found %d results\n", len(outputs))
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for i := 0; i < len(outputs); i += 3 {
		if i+2 < len(outputs) {
			fmt.Fprintf(
				writer,
				"%s\t%s\t%s\n",
				outputs[i],
				outputs[i+1],
				outputs[i+2],
			)
		} else if i+1 < len(outputs) {
			fmt.Fprintf(
				writer,
				"%s\t%s\t\n",
				outputs[i],
				outputs[i+1],
			)
		} else {
			fmt.Fprintf(
				writer,
				"%s\t\t\n",
				outputs[i],
			)
		}
	}
	writer.Flush()

	return nil
}

// Example modrinth url:
// https://api.modrinth.com/v2/search
// 		?limit=100
// 		&index=relevance
// 		&query=Carpet
// 		&facets=[
// 			[
// 				"categories:'forge'","categories:'fabric'","categories:'quilt'",
// 				"categories:'liteloader'","categories:'modloader'","categories:'rift'",
// 				"categories:'neoforge'"
// 			],
// 			["client_side:optional","client_side:unsupported"],
// 			["server_side:optional","server_side:required"],
// 			["project_type:mod"]
// 		]

// https://api.modrinth.com/v2/search?limit=20&index=relevance&query=Carpet&facets=%5B%5B%22categories:'forge'%22,%22categories:'fabric'%22,%22categories:'quilt'%22,%22categories:'liteloader'%22,%22categories:'modloader'%22,%22categories:'rift'%22,%22categories:'neoforge'%22%5D,%5B%22project_type:mod%22%5D%5D

func constructModrinthSearchUrl(q string) (url string) {
	var urlBuilder strings.Builder
	templateUrl, _ := template.New("template_url").Parse(`https://api.modrinth.com/v2/search?limit=100&index=relevance&query={{.}}&facets=%5B%5B%22categories:'forge'%22,%22categories:'fabric'%22,%22categories:'quilt'%22,%22categories:'liteloader'%22,%22categories:'modloader'%22,%22categories:'rift'%22,%22categories:'neoforge'%22%5D,%5B%22project_type:mod%22%5D%5D`)
	_ = templateUrl.Execute(&urlBuilder, q)
	return urlBuilder.String()
}

func searchModrinth(q string, loader string) (result types.ModrinthSearchRes) {
	switch loader {
	case "fabric":
		// fmt.Println("Fetching " + constructModrinthSearchUrl(q))
		resp, _ := http.Get(constructModrinthSearchUrl(q))
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &result)
		resp.Body.Close()
	case "forge":
		// TODO: Implement forge search
		fmt.Println("Not yet implemented")
	default:
		fmt.Println("Invalid loader")
	}
	return
}

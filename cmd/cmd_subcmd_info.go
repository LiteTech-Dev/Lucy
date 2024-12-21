package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"lucy/mcdr"
	"lucy/modrinth"
	"lucy/syntax"
	"lucy/types"
	"net/http"
	"os"
	"text/tabwriter"
)

var SubcmdInfo = &cli.Command{
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
	Action: ActionInfo,
}

func ActionInfo(ctx context.Context, cmd *cli.Command) error {
	// TODO: Error handling
	_, p := syntax.Parse(cmd.Args().First())

	switch p.Platform {
	case syntax.AllPlatform:
		// TODO: Wide range search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.PackageName))
		modrinthProject := &types.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, modrinthProject)
		generateInfoOutput(modrinthProject)
	case syntax.Fabric:
		// TODO: Fabric specific search
		res, _ := http.Get(modrinth.ConstructProjectUrl(p.PackageName))
		modrinthProject := &types.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, modrinthProject)
		generateInfoOutput(modrinthProject)
	case syntax.Forge:
		// TODO: Forge support
		println("Not yet implemented")
	case syntax.Mcdr:
		mcdrPlugin := mcdr.SearchMcdrPluginCatalogue(p.PackageName)
		if mcdrPlugin == nil {
			_ = fmt.Errorf("plugin not found")
			return nil
		}
		generateInfoOutput(mcdrPlugin)
	}

	return nil
}

func Bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func Magenta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func Faint(s string) string {
	return "\033[2m" + s + "\033[0m"
}

func printLabels(writer *tabwriter.Writer, labelTitle string, labels []string) {
	fmt.Fprintf(writer, "%s\t", Bold(Magenta(labelTitle)))
	lineLength := 0
	for _, label := range labels {
		fmt.Fprintf(writer, "%s", label)
		if label != labels[len(labels)-1] {
			fmt.Fprintf(writer, ", ")
		}
		lineLength += len(label) + 2
		if lineLength > 60 {
			fmt.Fprintf(writer, "\n")
			fmt.Fprintf(writer, "%s\t", Bold(Magenta("")))
			lineLength = 0
		}
	}
	if lineLength > 0 {
		fmt.Fprintf(writer, "\n")

	}
}

func generateInfoOutput(data interface{}) {
	writer := tabwriter.NewWriter(os.Stdout, 20, 4, 2, ' ', 0)

	switch v := data.(type) {
	case *types.ModrinthProject:
		fmt.Fprintf(writer, "%s\t%s\n", Bold(Magenta("Name")), v.Title)
		fmt.Fprintf(
			writer,
			"%s\t%s\n",
			Bold(Magenta("Description")),
			v.Description,
		)
		printLabels(writer, "Supported Versions", v.GameVersions)
		fmt.Fprintf(writer, "%s\t%s\n", Bold(Magenta("Source")), v.SourceUrl)
	case *types.McdrPluginInfo:
		fmt.Fprintf(writer, "%s\t%s\n", Bold(Magenta("Name")), v.Id)
		fmt.Fprintf(writer, "%s\t", Bold(Magenta("Author")))
		for i, author := range v.Authors {
			if i == 0 {
				fmt.Fprintf(writer, "%s ", author.Name)
				fmt.Fprintf(writer, "%s\n", Faint(author.Link))
			} else {
				fmt.Fprintf(writer, "%s\t", Bold(Magenta("")))
				fmt.Fprintf(writer, "%s ", author.Name)
				fmt.Fprintf(writer, "%s\n", Faint(author.Link))
			}
		}
		// fmt.Fprintf(writer, "%s\n", v.Authors[0].Name)
		// fmt.Fprintf(writer, "%s\t", Bold(Magenta("")))
		// fmt.Fprintf(writer, "%s\n", v.Authors[1].Name)
		fmt.Fprintf(writer, "%s\t%s\n", Bold(Magenta("Source")), v.Repository)
	default:
		fmt.Fprintf(writer, "Unsupported data type\n")
	}

	writer.Flush()
}

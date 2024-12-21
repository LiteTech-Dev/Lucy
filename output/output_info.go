package output

import (
	"fmt"
	"lucy/types"
	"os"
	"text/tabwriter"
)

func GenerateInfoOutput(data interface{}) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer writer.Flush()

	switch v := data.(type) {
	case *types.ModrinthProject:
		printInfo(writer, "Name", v.Title)
		printInfo(writer, "Description", v.Description)
		printInfo(writer, "Downloads", fmt.Sprintf("%d", v.Downloads))
		printKey(writer, "Game Versions")
		printLabels(writer, v.GameVersions, 60)
	case *types.McdrPluginInfo:
		printInfo(writer, "Name", v.Id)
		printKey(writer, "Authors")
		for i, author := range v.Authors {
			if i != 0 {
				printKey(writer, "")
			}
			printAuthors(writer, author.Name, author.Link)
		}
		printInfo(writer, "Source", v.Repository)
	default:
		panic("Invalid type")

	}
}

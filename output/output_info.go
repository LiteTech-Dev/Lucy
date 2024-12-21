package output

import (
	"fmt"
	"lucy/types"
)

func GenerateInfo(data interface{}) {
	defer keyValueWriter.Flush()

	switch v := data.(type) {
	case *types.ModrinthProject:
		printField("Name", v.Title)
		printField("Description", v.Description)
		printField("Downloads", fmt.Sprintf("%d", v.Downloads))
		printKey("Game Versions")
		printLabels(v.GameVersions, 60)
	case *types.McdrPluginInfo:
		printField("Name", v.Id)
		printKey("Authors")
		for i, author := range v.Authors {
			if i != 0 {
				printKey("")
			}
			printFieldWithAnnotation(author.Name, author.Link)
		}
		printField("Source", v.Repository)
	default:
		panic("Invalid type")
	}
}

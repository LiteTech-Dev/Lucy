package output

import (
	"fmt"
	"lucy/lucytypes"
)

func GenerateInfo(data interface{}) {
	defer keyValueWriter.Flush()

	switch v := data.(type) {
	case *lucytypes.ModrinthProject:
		printField("Name", v.Title)
		printField("Description", v.Description)
		printField("Downloads", fmt.Sprintf("%d", v.Downloads))
		printKey("Game Versions")
		printLabels(v.GameVersions, 60)
	case *lucytypes.McdrPluginInfo:
		printField("Name", v.Id)
		printKey("Authors")
		for i, author := range v.Authors {
			if i != 0 {
				printKey("")
			}
			printValueAnnot(author.Name, author.Link)
		}
		printField("Source", v.Repository)
	default:
		panic("Invalid type")
	}
}

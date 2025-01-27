package output

import (
	"fmt"
	"lucy/apitypes"
)

func GenerateInfo(data interface{}) {
	defer keyValueWriter.Flush()
	switch v := data.(type) {
	case *apitypes.ModrinthProject:
		SourceInfo("https://modrinth.com")
		ShortTextField("Name", v.Title)
		ShortTextField("Description", v.Description)
		ShortTextField("Downloads", fmt.Sprintf("%d", v.Downloads))
		LabelsField("Game Versions", v.GameVersions, 60)
	case *apitypes.McdrPluginInfo:
		SourceInfo("https://mcdreforged.com")
		ShortTextField("Name", v.Id)
		for i, author := range v.Authors {
			if i == 0 {
				ShortTextFieldWithAnnot("Author", author.Name, author.Link)
			} else {
				ShortTextFieldWithAnnot("", author.Name, author.Link)
			}
		}
		ShortTextField("Source", v.Repository)
	default:
		panic("Invalid data type")
	}
}

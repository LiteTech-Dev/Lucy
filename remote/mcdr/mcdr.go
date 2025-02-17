package mcdr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v50/github"
	"lucy/datatypes"
	"lucy/lucytypes"
	"lucy/syntax"
	"path"
)

func mcdrPluginInfoToPackageInfo(s *datatypes.McdrPluginInfo) *lucytypes.Package {
	name := lucytypes.PackageName(s.Id)

	info := &lucytypes.Package{
		Id: lucytypes.PackageId{
			Platform: lucytypes.Mcdr,
			Name:     name,
			Version:  lucytypes.LatestVersion,
		},
	}

	return info
}

func SearchMcdrPluginCatalogue(search lucytypes.PackageName) (
	pluginInfo *datatypes.McdrPluginInfo,
	err error,
) {
	plugins := getMcdrPluginCatalogue()

	for _, plugin := range plugins {
		p := syntax.Parse(*plugin.Name)
		if p.Name == search {
			return getMcdrPluginInfo(*plugin.Path), nil
		}
	}

	return nil, fmt.Errorf("plugin not found")
}

func getMcdrPluginCatalogue() []*github.RepositoryContent {
	ctx := context.Background()
	client := github.NewClient(nil)

	_, directoryContent, _, _ := client.Repositories.GetContents(
		ctx,
		"MCDReforged",
		"PluginCatalogue",
		"plugins",
		nil,
	)

	return directoryContent
}

func getMcdrPluginInfo(pluginPath string) (pluginInfo *datatypes.McdrPluginInfo) {
	ctx := context.Background()
	client := github.NewClient(nil)

	fileContent, _, _, _ := client.Repositories.GetContents(
		ctx,
		"MCDReforged",
		"PluginCatalogue",
		path.Join(pluginPath, "plugin_info.json"),
		nil,
	)
	pluginInfo = &datatypes.McdrPluginInfo{}
	content, _ := fileContent.GetContent()
	_ = json.Unmarshal([]byte(content), pluginInfo)

	return
}

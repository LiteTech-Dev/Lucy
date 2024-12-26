package mcdr

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v50/github"
	"lucy/lucytypes"
	"lucy/syntax"
	"path"
)

func SearchMcdrPluginCatalogue(slug syntax.PackageName) (pluginInfo *lucytypes.McdrPluginInfo) {
	plugins := getMcdrPluginCatalogue()

	for _, plugin := range plugins {
		_, p := syntax.Parse(*plugin.Name)
		if p.PackageName == slug {
			return getMcdrPluginInfo(*plugin.Path)
		}
	}

	return nil
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

func getMcdrPluginInfo(pluginPath string) (pluginInfo *lucytypes.McdrPluginInfo) {
	ctx := context.Background()
	client := github.NewClient(nil)

	fileContent, _, _, _ := client.Repositories.GetContents(
		ctx,
		"MCDReforged",
		"PluginCatalogue",
		path.Join(pluginPath, "plugin_info.json"),
		nil,
	)
	pluginInfo = &lucytypes.McdrPluginInfo{}
	content, _ := fileContent.GetContent()
	_ = json.Unmarshal([]byte(content), pluginInfo)

	return
}

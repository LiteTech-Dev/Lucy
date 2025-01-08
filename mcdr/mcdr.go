package mcdr

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v50/github"
	"lucy/apitypes"
	"lucy/syntax"
	"path"
)

func SearchMcdrPluginCatalogue(slug syntax.PackageName) (pluginInfo *apitypes.McdrPluginInfo) {
	plugins := getMcdrPluginCatalogue()

	for _, plugin := range plugins {
		p, _ := syntax.Parse(*plugin.Name)
		if p.Name == slug {
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

func getMcdrPluginInfo(pluginPath string) (pluginInfo *apitypes.McdrPluginInfo) {
	ctx := context.Background()
	client := github.NewClient(nil)

	fileContent, _, _, _ := client.Repositories.GetContents(
		ctx,
		"MCDReforged",
		"PluginCatalogue",
		path.Join(pluginPath, "plugin_info.json"),
		nil,
	)
	pluginInfo = &apitypes.McdrPluginInfo{}
	content, _ := fileContent.GetContent()
	_ = json.Unmarshal([]byte(content), pluginInfo)

	return
}

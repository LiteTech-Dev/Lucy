package output

import (
	"lucy/tools"
	"lucy/types"
	"strconv"
)

func GenerateStatus(data types.ServerInfo) {
	defer keyValueWriter.Flush()

	// Print game version
	printKey("Minecraft")
	printValueAnnot(data.Executable.GameVersion, data.Executable.Path)

	// Print active status
	if data.Activity != nil {
		printKey("Activity")
		printValueAnnot(
			"Currently Running",
			"pid "+strconv.Itoa(data.Activity.Pid),
		)
	} else {
		printField("Activity", "Inactive")
	}

	// Print mod loader types and version
	printField("Modding", captalize(string(data.Executable.Type)))

	// Print MCDR status
	if data.Modules.Mcdr != nil {
		printField("MCDR", "Installed")
		printKey("MCDR Plugins")
		printLabels(data.Modules.Mcdr.PluginPaths, 1)
	}

	// Print lucy status
	printField(
		"Lucy",
		tools.Trenary(
			func() bool { return data.HasLucy },
			"Installed",
			"Not Installed",
		),
	)
}

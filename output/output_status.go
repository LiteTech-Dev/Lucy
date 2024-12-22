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
	if data.IsRunning {
		printKey("Activity")
		printValueAnnot("Currently Running", "pid "+strconv.Itoa(data.Pid))
	} else {
		printField("Activity", "Inactive")
	}

	// Print mod loader types and version
	printKey("Modding")
	printValueAnnot(
		captalize(data.Executable.ModLoaderType),
		tools.Trenary(
			func() bool { return data.Executable.ModLoaderType != "vanilla" },
			"v"+data.Executable.ModLoaderVersion,
			"",
		),
	)

	// Print MCDR status
	if data.HasMcdr {
		printField("MCDR", "Installed")
		printKey("MCDR Plugins")
		printLabels(data.McdrPluginPaths, 1)
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

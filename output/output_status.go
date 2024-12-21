package output

import (
	"lucy/tools"
	"lucy/types"
)

func GenerateStatus(data types.ServerInfo) {
	defer keyValueWriter.Flush()

	// Print game version
	printKey("Minecraft")
	printFieldWithAnnotation(data.Executable.GameVersion, data.Executable.Path)

	// Print active status
	printField(
		"Activity",
		tools.Trenary(
			func() bool { return data.IsRunning },
			"Active",
			"Inactive",
		),
	)

	// Print mod loader types and version
	printKey("Modding")
	printFieldWithAnnotation(
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

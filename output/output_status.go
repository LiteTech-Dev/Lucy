package output

import (
	"fmt"
	"lucy/lucytypes"
	"lucy/syntaxtypes"
	"lucy/tools"
)

func GenerateStatus(data *lucytypes.ServerInfo) {
	defer keyValueWriter.Flush()

	// Print game version
	printKey("Minecraft")
	printValueAnnot(data.Executable.GameVersion, data.Executable.Path)

	// Print active status
	if data.Activity != nil {
		printKey("Activity")
		printValueAnnot(
			"Currently Running",
			fmt.Sprintf("PID: %d", data.Activity.Pid),
		)
	} else {
		printField("Activity", "Inactive")
	}

	// Print mod loader types and version
	if data.Executable.Platform != syntaxtypes.Minecraft {
		printField("Modding", data.Executable.Platform.String())
	}

	// Print MCDR status
	if data.Mcdr != nil {
		printField("MCDR", "Installed")
		printKey("MCDR Plugins")
		printLabels(data.Mcdr.PluginPaths, 1)
	}

	// Print lucy status
	printField(
		"Lucy",
		tools.Ternary(
			func() bool { return data.HasLucy },
			"Installed",
			"Not Installed",
		),
	)
}

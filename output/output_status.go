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
	ShortTextFieldWithAnnot(
		syntaxtypes.Minecraft.String(),
		data.Executable.GameVersion,
		data.Executable.Path,
	)

	// Print active status
	if data.Activity != nil {
		ShortTextFieldWithAnnot(
			"Activity",
			"Currently Running",
			fmt.Sprintf("PID: %d", data.Activity.Pid),
		)
	} else {
		ShortTextField("Activity", "Inactive")
	}

	// Print mod loader types and version
	if data.Executable.Platform != syntaxtypes.Minecraft {
		ShortTextField("Modding", data.Executable.Platform.String())
	}

	// Print MCDR status
	if data.Mcdr != nil {
		ShortTextField("MCDR", "Installed")
		LabelsField("MCDR Plugins", data.Mcdr.PluginPaths, 1)
	}

	// Print lucy status
	ShortTextField(
		"Lucy",
		tools.Ternary(
			func() bool { return data.HasLucy },
			"Installed",
			"Not Installed",
		),
	)
}

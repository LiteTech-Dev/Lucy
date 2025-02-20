package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"lucy/local"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/tools"
)

var subcmdStatus = &cli.Command{
	Name:   "status",
	Usage:  "Display basic information of the current server",
	Action: tools.Decorate(actionStatus, globalFlagsDecorator),
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Output in raw JSON format",
		},
	},
}

var actionStatus cli.ActionFunc = func(
	_ context.Context,
	cmd *cli.Command,
) error {
	serverInfo := local.GetServerInfo()
	if cmd.Bool("debug") {
		tools.PrintAsJson(serverInfo)
	} else {
		output.Flush(serverInfoToStatus(&serverInfo))
	}
	return nil
}

// Output order:
// 1. Game Version
// 2. Executable Path
// 3. Activity Status
// 4. Modding Platform
// 5. Mods

const statusFieldCount = 5

const (
	statusFieldGameVersion = iota
	statusFieldExecutablePath
	statusFieldActivity
	statusFieldModdingPlatform
	statusFieldMods
)

func serverInfoToStatus(data *lucytypes.ServerInfo) *lucytypes.OutputData {
	status := &lucytypes.OutputData{
		Fields: make([]lucytypes.Field, statusFieldCount),
	}

	status.Fields[statusFieldGameVersion] = &output.FieldShortText{
		Title: "Game Version",
		Text:  data.Executable.GameVersion,
	}

	status.Fields[statusFieldExecutablePath] = &output.FieldShortText{
		Title: "Executable",
		Text:  data.Executable.Path,
	}

	if data.Activity != nil {
		status.Fields[statusFieldActivity] = &output.FieldAnnotatedShortText{
			Title: "Activity",
			Text: tools.Ternary(
				data.Activity.Active,
				"Active",
				"Inactive",
			),
			Annotation: tools.Ternary(
				data.Activity.Active,
				fmt.Sprintf("PID: %d", data.Activity.Pid),
				"",
			),
		}
	} else {
		status.Fields[statusFieldActivity] = &output.FieldShortText{
			Title: "Activity",
			Text:  tools.Dim("Unknown"),
		}
	}

	// Modding related fields only shown when modding platform detected
	if data.Executable.Platform != lucytypes.Minecraft {
		mods := make([]string, 0, len(data.Mods))
		modPaths := make([]string, 0, len(mods))
		if len(data.Mods) == 0 {
			mods = append(mods, tools.Dim("(None)"))
		}
		for _, mod := range data.Mods {
			mods = append(mods, mod.Id.FullString())
			modPaths = append(modPaths, mod.Local.Path)
		}

		status.Fields[statusFieldModdingPlatform] = &output.FieldShortText{
			Title: "Platform",
			Text:  data.Executable.Platform.Title(),
		}

		status.Fields[statusFieldMods] = &output.FieldMultiShortTextWithAnnot{
			Title:  "Mod List",
			Texts:  mods,
			Annots: modPaths,
		}
	} else {
		status.Fields[statusFieldModdingPlatform] = output.FieldNil
		status.Fields[statusFieldMods] = output.FieldNil
	}

	return status
}

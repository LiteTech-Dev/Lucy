package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/lucytypes"
	"lucy/output"
	"lucy/probe"
	"lucy/syntaxtypes"
	"lucy/tools"
)

var subcmdStatus = &cli.Command{
	Name:   "status",
	Usage:  "Display basic information of the current server",
	Action: actionStatus,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Output in raw JSON format",
		},
	},
}

func actionStatus(_ context.Context, cmd *cli.Command) error {
	serverInfo := probe.GetServerInfo()
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
		Title: "Executable Path",
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
		status.Fields[statusFieldActivity] = output.FieldNil
	}

	if data.Executable.Platform != syntaxtypes.Minecraft {
		mods := make([]string, 0, len(data.Mods))
		if len(data.Mods) == 0 {
			mods = append(mods, "None")
		}
		for _, mod := range data.Mods {
			mods = append(mods, string(mod.Base.Name))
		}

		status.Fields[statusFieldModdingPlatform] = &output.FieldShortText{
			Title: "Modding Platform",
			Text:  data.Executable.Platform.String(),
		}
		status.Fields[statusFieldMods] = &output.FieldLabels{
			Title:  "Mods",
			Labels: mods,
		}
	} else {
		status.Fields[statusFieldModdingPlatform] = output.FieldNil
	}

	return status
}

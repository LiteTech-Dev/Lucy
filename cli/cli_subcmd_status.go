package cli

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"lucy/probe"
	"os"
	"strings"
)

var SubcmdStatus = &cli.Command{
	Name:   "status",
	Usage:  "Display basic information of the current server",
	Action: ActionStatus,
}

func ActionStatus(ctx context.Context, cmd *cli.Command) error {
	serverInfo := probe.GetServerInfo()

	// Print game version
	fmt.Printf("Minecraft v%s\n", serverInfo.Executable.GameVersion)

	// Print MCDR status
	if serverInfo.HasMcdr {
		fmt.Printf("Managed by MCDR\n")
	}

	// Print mod loader type and version
	fmt.Printf(
		"%s%s ",
		strings.ToUpper(serverInfo.Executable.ModLoaderType[:1]),
		serverInfo.Executable.ModLoaderType[1:],
	)
	if serverInfo.Executable.ModLoaderType != "vanilla" {
		fmt.Printf("v%s\n", serverInfo.Executable.ModLoaderVersion)
	} else {
		fmt.Printf("\n")
	}

	// Print lucy status
	if _, err := os.Stat(".lucy"); err == nil {
		fmt.Println("Lucy installed")
	} else if os.IsNotExist(err) {
		fmt.Println("Lucy not installed")
		fmt.Println("Run `lucy init` to install Lucy")
	}

	return nil
}

package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"github.com/urfave/cli/v3"
	"io"
	"lucy/types"
	"net/http"
	"os"
	"text/tabwriter"
)

func SubcmdInfo(ctx context.Context, cmd *cli.Command) error {
	platform, packageName := parsePackageSyntax(cmd.Args().First())
	switch platform {
	case "":
		return errors.New("invalid query format")
	case "all":
		// TODO: Wide range search
		res, _ := http.Get(constructMorinthInfoURL(packageName))
		modrinthProject := types.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, &modrinthProject)
		generateInfoOutput(modrinthProject)
	case "fabric":
		// TODO: Fabric specific search
		res, _ := http.Get(constructMorinthInfoURL(packageName))
		modrinthProject := types.ModrinthProject{}
		data, _ := io.ReadAll(res.Body)
		json.Unmarshal(data, &modrinthProject)
		generateInfoOutput(modrinthProject)
	case "forge":
		println("Not yet implemented")
	}
	return nil
}

func generateInfoOutput(modrinthProject types.ModrinthProject) {
	writer := tabwriter.NewWriter(os.Stdout, 40, 4, 2, ' ', 0)
	const maxWidth = 80
	wrappedBody := wordwrap.WrapString(modrinthProject.Body, uint(maxWidth))
	fmt.Fprintln(writer, wrappedBody)
	writer.Flush()
}

func constructMorinthInfoURL(packageName string) (url string) {
	return "https://api.modrinth.com/v2/project/" + packageName
}

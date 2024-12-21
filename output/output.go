package output

import (
	"fmt"
	"lucy/types"
	"os"
	"text/tabwriter"
)

var keyValueWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

func printKey(title string) {
	fmt.Fprintf(keyValueWriter, "%s\t", bold(mangeta(title)))
}

func printValue(value string) {
	fmt.Fprintf(keyValueWriter, "%s\n", value)
}

func printFieldWithAnnotation(value string, annotation string) {
	fmt.Fprintf(keyValueWriter, "%s %s\n", value, faint(annotation))
}

func printField(key string, value string) {
	fmt.Fprintf(keyValueWriter, "%s\t%s\n", bold(mangeta(key)), value)
}

func printLabels(labels []string, maxWidth int) {
	if len(labels) == 0 {
		fmt.Fprintf(keyValueWriter, "\n")
	} else if len(labels) == 1 {
		printValue(labels[0])
		return
	}

	width := 0
	for _, label := range labels {
		fmt.Fprintf(keyValueWriter, "%s", label)
		if label != labels[len(labels)-1] {
			fmt.Fprintf(keyValueWriter, ", ")
		}
		width += len(label) + 2
		if width > maxWidth {
			fmt.Fprintf(keyValueWriter, "\n%s\t", bold(mangeta("")))
			width = 0
		}
	}
	if width > 0 {
		fmt.Fprintf(keyValueWriter, "\n")
	}
}

func printVersions(
	versions []string,
	versionType types.MinecraftVersion, maxWidth int,
) {
	// TODO: filter by version type
	printLabels(versions, maxWidth)
}

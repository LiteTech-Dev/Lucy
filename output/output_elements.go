package output

import (
	"fmt"
	"lucy/tools"
	"os"
	"text/tabwriter"
)

const debugOutput = false

var keyValueWriter = tabwriter.NewWriter(
	os.Stdout,
	0,
	0,
	2,
	' ',
	tools.Ternary(debugOutput, tabwriter.Debug, 0),
)

func key(title string) {
	fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Mangeta(title)))
}

func value(value string) {
	fmt.Fprintf(keyValueWriter, "%s", value)
}

func inlineAnnot(annotation string) {
	fmt.Fprintf(keyValueWriter, "\t%s", tools.Dim(annotation))
}

func annot(value string) {
	fmt.Fprintf(keyValueWriter, "%s", tools.Dim(value))
}

func newLine() {
	fmt.Fprintf(keyValueWriter, "\n")
}

func tab() {
	fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Mangeta("")))
}

func flush() {
	keyValueWriter.Flush()
}

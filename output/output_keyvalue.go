package output

import (
	"fmt"
	"lucy/tools"
	"os"
	"text/tabwriter"
)

var keyValueWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

func key(title string) {
	fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Mangeta(title)))
}

func value(value string) {
	fmt.Fprintf(keyValueWriter, "%s", value)
}

func annot(annotation string) {
	fmt.Fprintf(keyValueWriter, "\t%s", tools.Dim(annotation))
}

func dim(value string) {
	fmt.Fprintf(keyValueWriter, "%s", tools.Dim(value))
}

func newLine() {
	fmt.Fprintf(keyValueWriter, "\n")
}

func tab() {
	fmt.Fprintf(keyValueWriter, "\t")
}

package output

import (
	"fmt"
	"lucy/types"
	"text/tabwriter"
)

func printKey(writer *tabwriter.Writer, title string) {
	fmt.Fprintf(writer, "%s\t", Bold(Magenta(title)))
}

func printAuthors(writer *tabwriter.Writer, name string, link string) {
	fmt.Fprintf(writer, "%s %s\n", name, Faint(link))
}

func printInfo(writer *tabwriter.Writer, key string, value string) {
	fmt.Fprintf(writer, "%s\t%s\n", Bold(Magenta(key)), value)
}

func printLabels(writer *tabwriter.Writer, labels []string, maxWidth int) {
	width := 0
	for _, label := range labels {
		fmt.Fprintf(writer, "%s", label)
		if label != labels[len(labels)-1] {
			fmt.Fprintf(writer, ", ")
		}
		width += len(label) + 2
		if width > maxWidth {
			fmt.Fprintf(writer, "\n%s\t", Bold(Magenta("")))
			width = 0
		}
	}
	if width > 0 {
		fmt.Fprintf(writer, "\n")

	}
}


package output

import (
	"fmt"
	"lucy/tools"
	"os"
	"text/tabwriter"
)

var keyValueWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

func printTitle(title string) {
	fmt.Fprintf(keyValueWriter, "%s\t", tools.Bold(tools.Mangeta(title)))
}

func printText(value string) {
	fmt.Fprintf(keyValueWriter, "%s", value)
}

func printNewLine() {
	fmt.Fprintf(keyValueWriter, "\n")
}

func printAnnotation(annotation string) {
	fmt.Fprintf(keyValueWriter, "%s", tools.Dim(annotation))
}

// ShortTextFieldWithAnnot prints a key-value pair with an annotation
func ShortTextFieldWithAnnot(key string, text string, annotation string) {
	ShortTextField(key, text)
	printAnnotation(annotation)
	printNewLine()
}

// ShortTextField prints a key-value pair
func ShortTextField(key string, text string) {
	printTitle(key)
	printText(text)
	printNewLine()
}

// LabelsField prints a list of labels
func LabelsField(title string, labels []string, maxWidth int) {
	if len(labels) == 0 {
		return
	} else if len(labels) == 1 {
		printTitle(title)
		printText(labels[0])
		return
	}

	printTitle(title)
	width := 0
	for _, label := range labels {
		fmt.Fprintf(keyValueWriter, "%s", label)
		if label != labels[len(labels)-1] {
			fmt.Fprintf(keyValueWriter, ", ")
		}
		width += len(label) + 2
		if width > maxWidth {
			fmt.Fprintf(keyValueWriter, "\n%s\t", tools.Bold(tools.Mangeta("")))
			width = 0
		}
	}
	if width > 0 {
		fmt.Fprintf(keyValueWriter, "\n")
	}
}

// TODO: Implement this

func LongTextField(title string, text string, maxWidth int) {}

// TODO: Implement this

func printVersions(
title string,
versions []string,
maxWidth int,
showAll bool,
) {
	LabelsField(title, versions, maxWidth)
}

package output

import (
	"fmt"
	"lucy/syntaxtypes"
	"lucy/tools"
)

// ShortTextFieldWithAnnot prints a key-value pair with an annotation
func ShortTextFieldWithAnnot(title string, text string, annotation string) {
	key(title)
	value(text)
	annot(annotation)
	newLine()
}

// ShortTextField prints a key-value pair
func ShortTextField(title string, text string) {
	key(title)
	value(text)
	newLine()
}

// LabelsField prints a list of labels
func LabelsField(title string, labels []string, maxWidth int) {
	if len(labels) == 0 {
		return
	} else if len(labels) == 1 {
		key(title)
		value(labels[0])
		return
	}

	key(title)
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

func SourceInfo(source string) {
	dim("(Source: " + tools.Underline(source) + ")")
	newLine()
}

// TODO: Implement this

func LongTextField(title string, text string, maxWidth int) {}

// TODO: Implement this

func VersionsField(
title string,
versions []syntaxtypes.PackageVersion,
maxWidth int,
releaseOnly bool,
) {
}

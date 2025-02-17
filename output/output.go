// Package output is a key-value based commandline output framework. It uses multiple
// different types of field to generate output.
//
// Note the field will not show if its content is empty
package output

import (
	"lucy/lucytypes"
	"lucy/tools"
	"strings"
)

func SourceInfo(source string) {
	annot("(Source: " + tools.Underline(source) + ")")
	newLine()
}

// Separator prints a separator line. A length of 0 will print a line of 66%
// terminal width.
//
// Separator also adjusts itself so it does not exceed the terminal width.
//
// Use dim to control whether the separator is dimmed.
func Separator(len int, dim bool) {
	if len == 0 {
		len = tools.TermWidth() * 2 / 3
	} else if len > tools.TermWidth() {
		len = tools.TermWidth()
	}

	sep := strings.Repeat("-", len)
	if dim {
		annot(sep)
	} else {
		value(sep)
	}
	newLine()
}

type FieldShortText struct {
	Title string
	Text  string
}

func (f *FieldShortText) Output() {
	key(f.Title)
	value(f.Text)
	newLine()
}

type FieldAnnotatedShortText struct {
	Title      string
	Text       string
	Annotation string
}

func (f *FieldAnnotatedShortText) Output() {
	key(f.Title)
	value(f.Text)
	inlineAnnot(f.Annotation)
	newLine()
}

var FieldNil = &fieldNil{}

type fieldNil struct{}

func (f *fieldNil) Output() {}

// FieldLabels is a field that contains a title and a list of labels. If the
// maxWidth is 0, it defaults to max(33% of terminal width, 40)
type FieldLabels struct {
	Title    string
	Labels   []string
	MaxWidth int
}

func (f *FieldLabels) Output() {
	if len(f.Labels) == 0 {
		return
	}

	key(f.Title)
	if f.MaxWidth == 0 {
		f.MaxWidth = max(33*tools.TermWidth()/100, 40)
	}
	width := 0
	for i, label := range f.Labels {
		value(label)
		if i != len(f.Labels)-1 {
			value(", ")
		}
		width += len(label) + 2
		if width >= f.MaxWidth && i != len(f.Labels)-1 {
			newLine()
			tab()
			width = 0
		}
	}

	if width != 0 {
		newLine()
	}
}

type FieldDynamicColumnLabels struct {
	Title  string
	Labels []string
}

func (f *FieldDynamicColumnLabels) Output() {
	if len(f.Labels) == 0 {
		return
	}

	// This field should have a unique indent size so it's predictable. Therefore
	// we call flush() before output.
	flush()
	key(f.Title)

	maxLabelLen := 0
	for _, label := range f.Labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	columns := (tools.TermWidth() - 4) / (maxLabelLen + 2)
	if columns == 0 {
		columns = 1
	}

	for i, label := range f.Labels {
		value(label)
		if (i+1)%columns == 0 || i == len(f.Labels)-1 {
			newLine()
		}
		tab()
	}

	// After output, we call flush() to reset the indent size.
	flush()
}

// FieldMultiShortTextWithAnnot accepts 2 arrays, Texts and Annots. len(Texts) determines
// the length of the output. Any content in Annots after len(Texts) will be omitted.
type FieldMultiShortTextWithAnnot struct {
	Title  string
	Texts  []string
	Annots []string
}

func (f *FieldMultiShortTextWithAnnot) Output() {
	if len(f.Texts) == 0 {
		return
	}

	for i, t := range f.Texts {
		if i == 0 {
			key(f.Title)
		} else {
			tab()
		}
		value(t)
		if i < len(f.Annots) {
			inlineAnnot(f.Annots[i])
		}
		newLine()
	}
}

type FieldMultiShortText struct {
	Title string
	Texts []string
}

func (f *FieldMultiShortText) Output() {
	if len(f.Texts) == 0 {
		return
	}

	for i, t := range f.Texts {
		if i == 0 {
			key(f.Title)
		} else {
			tab()
		}
		value(t)
		newLine()
	}
}

// FieldCheckBox defaults to a red cross and green check when TrueText and
// FalseText is not specified.
type FieldCheckBox struct {
	Title     string
	Boolean   bool
	TrueText  string
	FalseText string
}

func (f *FieldCheckBox) Output() {
	key(f.Title)

	if f.TrueText == "" {
		f.TrueText = tools.Green("\u2713") // Check
	}
	if f.FalseText == "" {
		f.FalseText = tools.Red("\u2717") // X
	}

	if f.Boolean {
		value(f.TrueText)
	} else {
		value(f.FalseText)
	}
}

func Flush(data *lucytypes.OutputData) {
	for _, field := range data.Fields {
		field.Output()
	}
	flush()
}

package output

import (
	"fmt"
	"golang.org/x/term"
	"lucy/lucytypes"
	"lucy/tools"
)

func SourceInfo(source string) {
	annot("(Source: " + tools.Underline(source) + ")")
	newLine()
}

// Separator prints a separator line. A length of 0 will print a line of current
// terminal width
func Separator(length int) {
	if length == 0 {
		length, _, _ = term.GetSize(0)
	}
	for i := 0; i < length; i++ {
		fmt.Fprintf(keyValueWriter, "-")
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
	if len(f.Labels) == 1 {
		key(f.Title)
		value(f.Labels[0])
		newLine()
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
			tab()
		} else {
			tab()
		}
	}

	flush()
}

// TODO: Refactor to FieldMultiShortTextWithAnnot
type FieldPeople struct {
	Title  string
	People []struct {
		Name string
		Link string
	}
}

func (f *FieldPeople) Output() {
	if len(f.People) == 0 {
		return
	}
	if len(f.People) == 1 {
		key(f.Title)
		value(f.People[0].Name)
		inlineAnnot(tools.Underline(f.People[0].Link))
		newLine()
		return
	}

	for i, person := range f.People {
		if i == 0 {
			key(f.Title)
			value(person.Name)
			inlineAnnot(tools.Underline(f.People[i].Link))
		} else {
			tab()
			value(person.Name)
			inlineAnnot(tools.Underline(f.People[i].Link))
		}
		if i != len(f.People)-1 {
			newLine()
		}
	}
}

func GenerateOutput(data *lucytypes.OutputData) {
	for _, field := range data.Fields {
		field.Output()
	}
	flush()
}

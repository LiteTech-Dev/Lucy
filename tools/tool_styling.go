package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	StyleReset = iota
	StyleBold
	StyleDim
	StyleItalic
	StyleUnderline
	StyleBlackText = iota + 25
	StyleRedText
	StyleGreenText
	StyleYellowText
	StyleBlueText
	StyleMagentaText
	StyleCyanText
	StyleWhiteText
)

func styleFactory(i int) func(any) string {
	return func(v any) string {
		s := v.(string)
		return fmt.Sprintf("\u001B[%dm%s\u001B[%dm", i, s, StyleReset)
	}
}

func Capitalize(v interface{}) string {
	s := v.(string)
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

var (
	Bold      = styleFactory(StyleBold)
	Dim       = styleFactory(StyleDim)
	Italic    = styleFactory(StyleItalic)
	Underline = styleFactory(StyleUnderline)
	Red       = styleFactory(StyleRedText)
	Green     = styleFactory(StyleGreenText)
	Yellow    = styleFactory(StyleYellowText)
	Blue      = styleFactory(StyleBlueText)
	Mangeta   = styleFactory(StyleMagentaText)
	Cyan      = styleFactory(StyleCyanText)
)

// PrintAsJson is usually used for debugging purposes
func PrintAsJson(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

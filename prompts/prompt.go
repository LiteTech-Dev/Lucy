package prompts

import (
	"github.com/manifoldco/promptui"
	"lucy/types"
	"strings"
	"text/template"
)

var selectExecutableTemplate = &promptui.SelectTemplates{
	Active:   formatExecutableChoice("●", "bold"),
	Inactive: formatExecutableChoice("○", ""),
	Selected: formatExecutableChoice("✔︎", "bold"),
	Help:     "",
	FuncMap:  appendFuncMap(),
}

func formatExecutableChoice(bullet, style string) string {
	return `{{ "` + bullet + `" | blue }} {{ .Path | ` + style + ` }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .ModLoaderType "vanilla" }}Vanilla{{ else }}{{ .ModLoaderType | capitalize }} v{{ .ModLoaderVersion }}{{ end }})[0m`
}

func appendFuncMap() template.FuncMap {
	funcMap := promptui.FuncMap
	funcMap["capitalize"] = func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return funcMap
}

func PromptSelectExecutable(executables []*types.ServerExecutable) int {
	selectExecutable := promptui.Select{
		Label:     "Multiple executables detected, select one",
		Items:     executables,
		Templates: selectExecutableTemplate,
	}
	index, _, _ := selectExecutable.Run()
	PromptRememberExecutable()
	return index
}

func PromptRememberExecutable() bool {
	confirmRememberExecutable := promptui.Prompt{
		Label:     "Remember this executable",
		IsConfirm: true,
	}
	result, _ := confirmRememberExecutable.Run()
	return result == "true"
}

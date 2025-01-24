package output

import (
	"encoding/json"
	"github.com/manifoldco/promptui"
	"lucy/lucytypes"
	"text/template"
)

var selectExecutableTemplate = &promptui.SelectTemplates{
	Active:   `{{ "●" | blue }} {{ .Path | bold }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .ModLoaderType "vanilla" }}Vanilla{{ else }}{{ .ModLoaderType | capitalize }} v{{ .ModLoaderVersion }}{{ end }})[0m`,
	Inactive: `{{ "○" | blue }} {{ .Path }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .ModLoaderType "vanilla" }}Vanilla{{ else }}{{ .ModLoaderType | capitalize }} v{{ .ModLoaderVersion }}{{ end }})[0m`,
	Selected: `{{ "✔︎" | green }} {{ .Path | bold }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .ModLoaderType "vanilla" }}Vanilla{{ else }}{{ .ModLoaderType | capitalize }} v{{ .ModLoaderVersion }}{{ end }})[0m`,
	FuncMap: func() (funcMap template.FuncMap) {
		funcMap = promptui.FuncMap
		funcMap["capitalize"] = captalize
		return
	}(),
}

func PromptSelectExecutable(executables []lucytypes.ExecutableInfo) int {
	data, _ := json.MarshalIndent(executables[0], "", "  ")
	println(string(data))
	selectExecutable := promptui.Select{
		Label:     "Multiple valid executables detected, select one:",
		Items:     executables,
		Templates: selectExecutableTemplate,
	}
	index, _, _ := selectExecutable.Run()
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

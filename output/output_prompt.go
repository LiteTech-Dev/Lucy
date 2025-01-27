package output

import (
	"github.com/manifoldco/promptui"
	"lucy/lucytypes"
)

var selectExecutableTemplate = &promptui.SelectTemplates{
	Active:   `{{ "●" | blue }} {{ .Path | bold }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .Platform "minecraft" }}Vanilla{{ else }}{{ .Platform }} {{ .LoaderVersion }}{{ end }})[0m`,
	Inactive: `{{ "○" | blue }} {{ .Path }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .Platform "minecraft" }}Vanilla{{ else }}{{ .Platform }} {{ .LoaderVersion }}{{ end }})[0m`,
	Selected: `{{ "✔︎" | green }} {{ .Path | bold }} [2m(Minecraft {{ .GameVersion }}, {{ if eq .Platform "minecraft" }}vsanilla{{ else }}{{ .Platform }} {{ .LoaderVersion }}{{ end }})[0m`,
}

func PromptSelectExecutable(executables []*lucytypes.ExecutableInfo) int {
	// data, _ := json.MarshalIndent(executables, "", "  ")
	// println(string(data))
	selectExecutable := promptui.Select{
		Label:     "Multiple possible executables detected, select one",
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

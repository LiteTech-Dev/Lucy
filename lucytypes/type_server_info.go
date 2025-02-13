package lucytypes

import (
	"os/exec"
)

// ServerInfo components that do not exist, use an empty string. Note Executable
// must exist, otherwise the program will exit; therefore, it is not a pointer.
type ServerInfo struct {
	WorkPath    string
	SavePath    string
	ModPath     string
	Mods        []*Package
	HasLucy     bool
	Mcdr        *McdrInstallation
	Executable  *ExecutableInfo
	BootCommand *exec.Cmd
	Activity    *Activity
}

type ExecutableInfo struct {
	Path          string
	GameVersion   string
	Platform      Platform
	LoaderVersion string
	BootCommand   *exec.Cmd
}

type Activity struct {
	Active bool
	Pid    int
}

type McdrInstallation struct {
	PluginPaths []string
	PluginList  []*Package // TODO: Implement probe func
}

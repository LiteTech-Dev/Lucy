package lucytypes

import (
	"lucy/syntaxtypes"
	"os/exec"
)

// ServerInfo components that do not exist, use an empty string. Note Executable
// must exist, otherwise the program will exit; therefore, it is not a pointer.
type ServerInfo struct {
	WorkPath    string
	SavePath    string
	ModPath     string
	Mods        []*PackageInfo
	HasLucy     bool
	Mcdr        *Mcdr
	Executable  *ExecutableInfo
	BootCommand *exec.Cmd
	Activity    *Activity
}

type ExecutableInfo struct {
	Path          string
	GameVersion   string
	Platform      syntaxtypes.Platform
	LoaderVersion string
	BootCommand   *exec.Cmd
}

type Activity struct {
	Active bool
	Pid    int
}

type Mcdr struct {
	PluginPaths []string
}

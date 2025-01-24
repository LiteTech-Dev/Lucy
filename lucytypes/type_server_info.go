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
	HasLucy     bool
	Executable  ExecutableInfo
	BootCommand exec.Cmd
	Activity    *Activity
	Modules     *ServerModules
}

type ExecutableInfo struct {
	Path          string
	GameVersion   string
	Type          syntaxtypes.Platform
	LoaderVersion string
	BootCommand   *exec.Cmd
}

type Activity struct {
	Active bool
	Pid    int
}

type ServerModules struct {
	Mcdr   *Mcdr
	Fabric *Fabric
	Forge  *Forge
}

type Mcdr struct {
	Name        syntaxtypes.Platform
	PluginPaths []string
}

type Fabric struct {
	Name    syntaxtypes.Platform
	Version string
}

type Forge struct {
	Name    syntaxtypes.Platform
	Version string
}

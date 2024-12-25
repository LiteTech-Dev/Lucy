package types

import (
	"lucy/syntax"
	"os/exec"
)

// ServerInfo components that do not exist, use an empty string. Note Executable
// must exist, otherwise the program will exit; therefore, it is not a pointer.
type ServerInfo struct {
	ServerWorkPath string
	SavePath       string
	ModPath        string
	HasLucy        bool
	Executable     ServerExecutable
	BootCommand    *exec.Cmd
	Activity       *Activity
	Modules        *ServerModules
}

type ServerExecutable struct {
	Path        string
	GameVersion string
	BootCommand exec.Cmd
	Type        syntax.Platform
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
	Name        syntax.Platform
	PluginPaths []string
}

type Fabric struct {
	Name    syntax.Platform
	Version string
}

type Forge struct {
	Name    syntax.Platform
	Version string
}

package types

// ServerInfo components that do not exist, use an empty string
type ServerInfo struct {
	HasMcdr         bool
	McdrConfigPath  string
	McdrPluginPaths []string

	ServerWorkPath string
	ModPath        string

	Executable *ServerExecutable
}

type ServerExecutable struct {
	Path             string
	GameVersion      string
	ModLoaderType    string
	ModLoaderVersion string
}

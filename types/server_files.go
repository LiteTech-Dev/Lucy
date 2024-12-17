package types

// ServerInfo components that do not exist, use an empty string
type ServerInfo struct {
	HasMcdr        bool
	McdrConfigPath string

	ServerWorkPath  string
	ModPath         string
	McdrPluginPaths []string
	ExecutablePath  string

	GameVersion      string
	ModLoaderType    string
	ModLoaderVersion string
}

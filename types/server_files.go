package types

// ServerFiles components that do not exist, use an empty string
type ServerFiles struct {
	HasMcdr        bool
	McdrConfigPath string

	ServerWorkPath  string
	ModPath         string
	McdrPluginPaths []string
	ExecutablePath  string

	MinecraftVersion string
	ModLoaderType    string
	ModLoaderVersion string
}

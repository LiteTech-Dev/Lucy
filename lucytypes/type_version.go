package lucytypes

type MinecraftVersion string

const (
	MinecraftSnapshot         MinecraftVersion = "snapshot"
	MinecraftPre              MinecraftVersion = "pre"
	MinecraftReleaseCandidate MinecraftVersion = "rc"
	MinecraftRelease          MinecraftVersion = "release"
)

func ParseMinecraftVersion(version string) MinecraftVersion {
	return ""
}

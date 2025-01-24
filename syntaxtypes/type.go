package syntaxtypes

// Platform is an enum of several string constants. All platform is a package under
// itself, for example, "fabric/fabric" is a valid package, and is equivalent to
// "fabric". This literal is typically used when installing/upgrading a platform
// itself.
type Platform string

const (
	Minecraft   Platform = "minecraft"
	Fabric      Platform = "fabric"
	Forge       Platform = "forge"
	Neoforge    Platform = "neoforge"
	Mcdr        Platform = "mcdr"
	AllPlatform Platform = "all"
)

var Platforms = []Platform{
	Minecraft, Fabric, Forge, Neoforge, Mcdr, AllPlatform,
}

// PackageName is the slug of the package, using hyphens as separators. For example,
// "fabric-api". It is not case-sensitive, however lowercase is recommended. Underline
// '_' is equivalent to hyphen. The slug from a source API is preferred, if available.
// Otherwise, the slug is obtained from the executable file. No exceptions since
// a package must either exist on a remote API or user's local files. All Minecraft
// versions are valid package names. This literal is typically used when migrating
// to another Minecraft version.
type PackageName string

type Package struct {
	Platform Platform
	Name     PackageName
	Version  PackageVersion
}

// IsValid should be edited if you added a new platform.
func (p Platform) IsValid() bool {
	for _, valid := range Platforms {
		if p == valid {
			return true
		}
	}
	return false
}

// PackageVersion is the version of the package. If not specified, it defaults to
// "all". Here we expect mods and plugins to use semver (which they should). A
// known exception is Minecraft versions snapshots.
type PackageVersion string

var AllVersion PackageVersion = "all"

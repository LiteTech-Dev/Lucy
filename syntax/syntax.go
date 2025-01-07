// Package syntax defines the syntax for specifying packages and platforms.
//
// A package can either be specified by a string in the format of
// "platform/name@version". Only the name is required, both platform and version
// can be omitted.
//
// Valid Examples:
//   - carpet
//   - mcdr/prime-backup
//   - fabric/jade@1.0.0
//   - fabric@12.0
//   - minecraft@1.19 (recommended)
//   - minecraft/minecraft@1.16.5 (= minecraft@1.16.5)
//   - minecraft/1.14.3 (= minecraft@1.14.3)
//   - 1.8.9 (= minecraft@1.8.9)
package syntax

import (
	"golang.org/x/mod/semver"
	"lucy/lucyerrors"
	"strings"
)

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

// Validate should be edited if you added a new platform.
func (p Platform) Validate() bool {
	switch p {
	case Fabric, Forge, Neoforge, Mcdr, Minecraft, AllPlatform:
		return true
	default:
		return false
	}
}

// PackageVersion is the version of the package. If not specified, it defaults to
// "all". Most mods should use semver. An exception is Minecraft versions snapshots.
// Therefore, the type MinecraftVersion is defined.
type PackageVersion string

const AllVersion = PackageVersion("all")

type MinecraftVersion PackageVersion

func (v MinecraftVersion) IsSnapshot() bool {
	return semver.IsValid(string("v" + v))
}

// sanitize tolerates some common interchangeability between characters. This
// includes underscores, chinese full stops, and backslashes. It also converts
// uppercase characters to lowercase.
func sanitize(s string) (clean string) {
	clean = ""
	for _, char := range s {
		if char == '_' {
			clean += string('-')
		} else if char == '\\' {
			clean += string('/')
		} else if char == '。' {
			clean += string('.')
		} else if char >= 'A' && char <= 'Z' {
			clean += strings.ToLower(string(char))
		} else {
			clean += string(char)
		}
	}
	return
}

func Parse(s string) (err error, p *Package) {
	s = sanitize(s)
	slashSplit := strings.Split(s, "/")
	p = &Package{}
	var atSplit []string

	switch len(slashSplit) {
	case 0:
		return lucyerrors.EmptyPackageSyntaxError, nil
	case 1:
		p.Platform = AllPlatform
		atSplit = strings.Split(slashSplit[0], "@")
	case 2:
		p.Platform = Platform(slashSplit[0])
		if !p.Platform.Validate() {
			return lucyerrors.InvalidPlatformError, nil
		}
		atSplit = strings.Split(slashSplit[1], "@")
	default:
		return lucyerrors.PackageSyntaxError, nil
	}

	switch len(atSplit) {
	case 1:
		p.Name = PackageName(atSplit[0])
		p.Version = AllVersion
	case 2:
		p.Name = PackageName(atSplit[0])
		p.Version = PackageVersion(atSplit[1])
	default:
		return lucyerrors.PackageSyntaxError, nil
	}

	return
}

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
//   - 1.8.9 (= minecraft@1.8.9)
package syntax

import (
	"errors"
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

var (
	ESyntax = errors.New("invalid syntax")
)

func Parse(s string) (p *Package, err error) {
	s = sanitize(s)
	p.Platform, p.Name, p.Version, err = parseAt(s)
	if err != nil {
		return nil, err
	}
	return
}

// parseAt is called first since '@' operator always occur after '/' (equivalent
// to a lower priority)
func parseAt(s string) (
	pl Platform,
	n PackageName,
	v PackageVersion,
	err error,
) {
	split := strings.Split(s, "@")

	pl, n, err = parseSlash(split[0])
	if err != nil {
		return "", "", "", ESyntax
	}

	if len(split) == 1 {
		v = AllVersion
	} else if len(split) == 2 {
		v = PackageVersion(split[1])
	} else {
		return "", "", "", ESyntax
	}

	return
}

func parseSlash(s string) (pl Platform, n PackageName, err error) {
	split := strings.Split(s, "/")

	if len(split) == 1 {
		pl = AllPlatform
		n = PackageName(split[0])
		for _, platform := range Platforms {
			if PackageName(platform) == n {
				pl = platform
			}
		}
	} else if len(split) == 2 {
		pl = Platform(split[0])
		n = PackageName(split[1])
	} else {
		return "", "", ESyntax
	}

	return
}

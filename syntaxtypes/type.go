// Package syntaxtypes contains syntax-related constants and types. This package
// is created to avoid cyclic dependencies by separating syntax-related types and
// syntax parsing functions. As types are more fundamental, and were used in more
// packages, they are placed in this independent package. This package should not
// import any other packages in lucy.
package syntaxtypes

import (
	"fmt"
	"lucy/tools"
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

func (p Platform) String() string {
	return tools.Capitalize(p)
}

func (p Platform) IsAll() bool {
	return p == AllPlatform
}

// Valid should be edited if you added a new platform.
func (p Platform) Valid() bool {
	for _, valid := range Platforms {
		if p == valid {
			return true
		}
	}
	return false
}

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

func (p *Package) String() string {
	return fmt.Sprintln(
		tools.Ternary(
			func() bool { return p.Platform == AllPlatform },
			"", string(p.Platform)+"/",
		),
		string(p.Name),
		tools.Ternary(
			func() bool { return p.Version == LatestVersion || p.Version == AllVersion || p.Version == NoVersion },
			"", "@"+string(p.Version),
		),
	)
}

// PackageVersion is the version of a package. Here we expect mods and plugins
// use semver (which they should). A known exception is Minecraft snapshots.
type PackageVersion string

var AllVersion PackageVersion = "all"
var NoVersion PackageVersion = "none"
var LatestVersion PackageVersion = "latest"

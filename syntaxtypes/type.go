// Package syntaxtypes contains syntax-related constants and types. This package
// is created to avoid cyclic dependencies by separating syntax-related types and
// syntax parsing functions. As types are more fundamental, and were used in more
// packages, they are placed in this independent package. This package should not
// import any other packages in lucy.
package syntaxtypes

import (
	"fmt"
	"lucy/tools"
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

// String by default returns a capitalized string
func (p Platform) String() string {
	if p.IsAll() {
		return "Any"
	}
	if p.Valid() {
		return strings.ToUpper(string(p)[0:1]) + string(p)[1:]
	}
	return "Invalid Platform"
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
// "fabric-api".
//
// It is non-case-sensitive, though lowercase is recommended. Underlines '_' are
// equivalent to hyphens.
//
// A slug from a upstream API is preferred, if possible. Otherwise, the slug is
// obtained from the executable file. No exceptions since a package must either
// exist on a remote API or user's local files.
type PackageName string

type Package struct {
	Platform Platform
	Name     PackageName
	Version  PackageVersion
}

func (p *Package) String() string {
	return fmt.Sprintln(
		tools.TernaryFunc(
			func() bool { return p.Platform == AllPlatform },
			"", string(p.Platform)+"/",
		),
		string(p.Name),
		tools.TernaryFunc(
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

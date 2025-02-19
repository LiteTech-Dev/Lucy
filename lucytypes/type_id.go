package lucytypes

import (
	"fmt"
	"strings"

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

func (p Platform) Title() string {
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

// Title Replaces underlines or hyphens with spaces, then capitalize the first
// letter.
func (p PackageName) Title() string {
	return tools.Capitalize(strings.ReplaceAll(string(p), "-", " "))
}

func (p PackageName) String() string {
	return string(p)
}

type PackageId struct {
	Platform Platform
	Name     PackageName
	Version  PackageVersion
}

func (p *PackageId) NewPackage() *Package {
	return &Package{
		Id: PackageId{
			Platform: p.Platform,
			Name:     p.Name,
			Version:  p.Version,
		},
	}
}

func (p *PackageId) String() string {
	return fmt.Sprintln(
		tools.Ternary(
			p.Platform == AllPlatform,
			"", string(p.Platform)+"/",
		),
		string(p.Name),
		tools.Ternary(
			p.Version == LatestVersion || p.Version == AllVersion || p.Version == NoVersion,
			"", "@"+string(p.Version),
		),
	)
}

// PackageVersion is the version of a package. Here we expect mods and plugins
// use semver (which they should). A known exception is Minecraft snapshots.
type PackageVersion string

func (p PackageVersion) String() string {
	if p == AllVersion || p == "" {
		return "Any"
	}
	if p == NoVersion {
		return "None"
	}
	if p == LatestVersion {
		return "Latest"
	}
	return string(p)
}

var (
	AllVersion    PackageVersion = "all"
	NoVersion     PackageVersion = "none"
	LatestVersion PackageVersion = "latest"
)

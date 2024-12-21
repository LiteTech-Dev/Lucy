package syntax

// Package syntax:
// A package can either be specified by its name or string in the format of
// "platform/package@version". The package name is the slug of the package.
// There is a special syntax for platforms only used when updating them (migration):
// "platform@version", you must not specify any packages in this way. Minecraft
// is both a package and a platform.
// Example: carpet
// Example: fabric/carpet@1.0.0
// Example: mcdr/prime-backup
// Example: fabric@12.0
// Example: minecraft@1.19
// Example: minecraft/minecraft@1.16.5

import (
	"errors"
	"strings"
)

type Platform string
type PackageName string
type PackageVersion string
type Package struct {
	Platform
	PackageName
	PackageVersion
}

const (
	Minecraft   Platform = "minecraft"
	Fabric      Platform = "fabric"
	Forge       Platform = "forge"
	Neoforge    Platform = "neoforge"
	Mcdr        Platform = "mcdr"
	AllPlatform Platform = "all"
)

const (
	AllVersion PackageVersion = "all"
)

const (
	MinecraftAsPackage = PackageName(Minecraft)
	FabricAsPackage    = PackageName(Fabric)
	ForgeAsPackage     = PackageName(Forge)
	NeoforgeAsPackage  = PackageName(Neoforge)
)

var InvalidPlatformError = errors.New("invalid platform")
var PackageSyntaxError = errors.New("invalid package syntax")
var EmptyPackageSyntaxError = errors.New("empty package string")

// validatePlatform should be edited if you added anything
func validatePlatform(value string) error {
	switch Platform(value) {
	case Fabric, Forge, Neoforge, Mcdr, Minecraft, AllPlatform:
		return nil
	default:
		return InvalidPlatformError
	}
}

// Sanitize tolerates some common interchangeability between characters. This
// includes underscores, chinese full stops, and backslashes. It also converts
// uppercase characters to lowercase.
func Sanitize(str string) (cleanStr string) {
	cleanStr = ""
	for _, char := range str {
		if char == '_' {
			cleanStr += string('-')
		} else if char == '\\' {
			cleanStr += string('/')
		} else if char == '。' {
			cleanStr += string('.')
		} else if char >= 'A' && char <= 'Z' {
			cleanStr += strings.ToLower(string(char))
		} else {
			cleanStr += string(char)
		}
	}
	return
}

func Parse(str string) (err error, p *Package) {
	str = Sanitize(str)
	slashSplit := strings.Split(str, "/")
	p = &Package{}
	var atSplit []string

	switch len(slashSplit) {
	case 0:
		return EmptyPackageSyntaxError, nil
	case 1:
		p.Platform = AllPlatform
		atSplit = strings.Split(slashSplit[0], "@")
	case 2:
		if err := validatePlatform(slashSplit[0]); err != nil {
			return err, nil
		}
		p.Platform = Platform(slashSplit[0])
		atSplit = strings.Split(slashSplit[1], "@")
	default:
		return PackageSyntaxError, nil
	}

	switch len(atSplit) {
	case 1:
		p.PackageName = PackageName(atSplit[0])
		p.PackageVersion = AllVersion
	case 2:
		p.PackageName = PackageName(atSplit[0])
		p.PackageVersion = PackageVersion(atSplit[1])
	default:
		return PackageSyntaxError, nil
	}

	return
}

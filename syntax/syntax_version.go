package syntax

import (
	"encoding/json"
	"golang.org/x/mod/semver"
	"io"
	"lucy/apitypes"
	"net/http"
)

// PackageVersion is the version of the package. If not specified, it defaults to
// "all". Most mods should use semver. An exception is Minecraft versions snapshots.
// Therefore, the type MinecraftVersion is defined.
type PackageVersion interface {
	IsValidVersion() bool
}

type typeAllVersion string

func (s typeAllVersion) IsValidVersion() bool {
	return s == "all"
}

const AllVersion = typeAllVersion("all")

type MinecraftVersion string

func (s MinecraftVersion) IsValidVersion() bool {
	manifest, _ := GetVersionManifest()
	for _, v := range manifest.Versions {
		if string(s) == v.Id {
			return true
		}
	}
	return false
}

type ModVersion string

func (s ModVersion) IsValidVersion() bool {
	return semver.IsValid("v" + string(s))
}

const VersionManifestURL = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

// TODO: Remove the err return value
// TODO: Use tools.Memoize to cache the result

func GetVersionManifest() (manifest *apitypes.VersionManifest, err error) {
	manifest = &apitypes.VersionManifest{}

	// TODO: Add cache mechanism if http call is too slow or fails
	resp, err := http.Get(VersionManifestURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// CompareMinecraftVersions gives -1 when v1 is older than v2, 0 when they are
// the same, and 1 when v1 is newer than v2
func CompareMinecraftVersions(v1, v2 PackageVersion) int8 {
	if v1 == v2 {
		return 0
	}

	manifest, _ := GetVersionManifest()
	i1, i2 := -1, -1
	for i, v := range manifest.Versions {
		if v1 == MinecraftVersion(v.Id) {
			i1 = i
		}
		if v2 == MinecraftVersion(v.Id) {
			i2 = i
		}
	}

	if i1 == i2 {
		return 0
	}
	if i1 < i2 {
		return 1
	}
	if i1 > i2 {
		return -1
	}
	panic("wrong version")
}

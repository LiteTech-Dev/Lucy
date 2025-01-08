package syntax

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/mod/semver"
	"io"
	"lucy/apitypes"
	"net/http"
)

type VersionFormat string

const (
	MinecraftVersion VersionFormat = "minecraft"
	SemanticVersion  VersionFormat = "semver"
	tAllVersion      VersionFormat = "all"
)

var AllVersion = PackageVersion{tAllVersion, ""}

// PackageVersion is the version of the package. If not specified, it defaults to
// "all". Most mods should use semver. An exception is Minecraft versions snapshots.
// Therefore, the type MinecraftVersion is defined.
type PackageVersion struct {
	format VersionFormat
	raw    string
}

const VersionManifestURL = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

var (
	EInvalidVersionComparison = errors.New("invalid version comparison")
	EVersionNotFound          = errors.New("version do not exist")
)

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
// the same (or an error occurred), and 1 when v1 is newer than v2. 0 is returned
// when either v1 or v2 is AllVersion
func CompareMinecraftVersions(v1, v2 PackageVersion) (c int8, err error) {
	if v1.format == tAllVersion || v2.format == tAllVersion {
		return 0, nil
	}

	if v1.format != v2.format {
		return 0, fmt.Errorf(
			"%w between type %s and %s",
			EInvalidVersionComparison,
			v1.format,
			v2.format,
		)
	}

	if v1.raw == v2.raw {
		return 0, nil
	}

	if v1.format == MinecraftVersion {
		return compareMinecraftVersions(v1, v2)
	} else if v1.format == SemanticVersion {
		return int8(semver.Compare("v"+v1.raw, "v"+v2.raw)), nil
	}

	return 0, nil // unreachable
}

func compareMinecraftVersions(v1, v2 PackageVersion) (c int8, err error) {
	manifest, err := GetVersionManifest()
	if err != nil {
		return 0, err
	}

	i1, i2 := -1, -1
	for i, v := range manifest.Versions {
		if v1.raw == (v.Id) {
			i1 = i
		}
		if v2.raw == (v.Id) {
			i2 = i
		}
	}

	if i1 == i2 {
		return 0, nil
	}
	if i1 < i2 {
		return 1, nil
	}
	if i1 > i2 {
		return -1, nil
	}

	return 0, EInvalidVersionComparison
}

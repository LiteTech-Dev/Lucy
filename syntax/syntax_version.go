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

// PackageVersion is the version of the package. If not specified, it defaults to
// "all". Here we expect mods and plugins to use semver (which they should). A
// known exception is Minecraft versions snapshots.
type PackageVersion string

var AllVersion PackageVersion = "all"
var (
	EInvalidVersionComparison = errors.New("invalid version comparison")
	EVersionNotFound          = errors.New("version do not exist")
)

// TODO: Remove the err return value
// TODO: Use tools.Memoize to cache the result

const VersionManifestURL = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

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

// ComparePackageVersions gives -1 when v1 is older than v2, 0 when they are
// the same (or an error occurred), and 1 when v1 is newer than v2. 0 is returned
// when either v1 or v2 is AllVersion
func ComparePackageVersions(p1, p2 *Package) (c int8, err error) {
	v1, v2 := p1.Version, p2.Version

	if v1 == AllVersion || v2 == AllVersion {
		return 0, nil
	}

	if p1.Platform != p2.Platform {
		return 0, fmt.Errorf(
			"%w between package on platform %s and %s",
			EInvalidVersionComparison,
			p1.Platform,
			p2.Platform,
		)
	}

	if v1 == v2 {
		return 0, nil
	}

	if p1.Platform == Minecraft {
		return compareMinecraftVersions(v1, v2)
	}
	return int8(semver.Compare("v"+string(v1), "v"+string(v2))), nil
}

func compareMinecraftVersions(v1, v2 PackageVersion) (c int8, err error) {
	manifest, err := GetVersionManifest()
	if err != nil {
		return 0, err
	}

	i1, i2 := -1, -1
	for i, v := range manifest.Versions {
		if v1 == PackageVersion(v.Id) {
			i1 = i
		}
		if v2 == PackageVersion(v.Id) {
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

package mcdr

import (
	"lucy/apitypes"
	"lucy/lucytypes"
	"lucy/syntaxtypes"
)

func mcdrPluginInfoToPackageInfo(s *apitypes.McdrPluginInfo) *lucytypes.PackageInfo {
	name := syntaxtypes.PackageName(s.Id)

	info := &lucytypes.PackageInfo{
		Id: syntaxtypes.Package{
			Platform: syntaxtypes.Mcdr,
			Name:     name,
			Version:  syntaxtypes.LatestVersion,
		},
		Path:               "",    // Wait for plugin list detection
		Installed:          false, // Wait for plugin list detection
		Urls:               []lucytypes.PackageUrl{},
		Name:               name.String(),
		Description:        s.Introduction.EnUs,
		SupportedVersions:  []syntaxtypes.PackageVersion{syntaxtypes.AllVersion},
		SupportedPlatforms: []syntaxtypes.Platform{syntaxtypes.Mcdr},
	}

	return info
}

package mcdr

import (
	"lucy/apitypes"
	"lucy/lucytypes"
)

func mcdrPluginInfoToPackageInfo(s *apitypes.McdrPluginInfo) *lucytypes.Package {
	name := lucytypes.PackageName(s.Id)

	info := &lucytypes.Package{
		Id: lucytypes.PackageId{
			Platform: lucytypes.McdrInstallation,
			Name:     name,
			Version:  lucytypes.LatestVersion,
		},
		Path:               "",    // Wait for plugin list detection
		Installed:          false, // Wait for plugin list detection
		Urls:               []lucytypes.PackageUrl{},
		Name:               name.String(),
		Description:        s.Introduction.EnUs,
		SupportedVersions:  []lucytypes.PackageVersion{lucytypes.AllVersion},
		SupportedPlatforms: []lucytypes.Platform{lucytypes.McdrInstallation},
	}

	return info
}

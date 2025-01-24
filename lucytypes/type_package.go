package lucytypes

import "lucy/syntaxtypes"

// PackageInfo and syntaxtypes.Package's difference is that PackageInfo
// only exists for packages already installed, while Package is an identifier
// for all packages. Therefore, here we define it as a super set of Package.
type PackageInfo struct {
	Base              syntaxtypes.Package
	Path              string
	SupportedVersions []syntaxtypes.PackageVersion
}

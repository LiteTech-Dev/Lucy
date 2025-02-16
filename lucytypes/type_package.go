package lucytypes

type PackageUrlType uint8

const (
	FileUrl PackageUrlType = iota
	HomepageUrl
	SourceUrl
	WikiUrl
	OthersUrl
)

func (p PackageUrlType) String() string {
	switch p {
	case FileUrl:
		return "File"
	case HomepageUrl:
		return "Homepage"
	case SourceUrl:
		return "Source"
	case WikiUrl:
		return "Wiki"
	case OthersUrl:
		return "URL"
	default:
		return "Unknown"
	}
}

type PackageUrl struct {
	Name string
	Type PackageUrlType
	Url  string
}

// Package is a package identifier with its related information.
//
// In principle, only package remote and local can provide a package.
type Package struct {
	// Id is the basic package identifier
	Id PackageId

	// Optional attributions
	Information  *PackageInformation
	Dependencies *PackageDependencies
	Local        *PackageInstallation
	Remote       *PackageRemote
}

// PackageDependencies is one of the optional attributions that can be added to
// a Package struct. It is usually used in any command that requires operating
// local packages, such as `lucy install` or `lucy remove`.
type PackageDependencies struct {
	SupportedVersions   []PackageVersion
	SupportedPlatforms  []Platform
	PackageDependencies []PackageId
}

// PackageInformation is a struct that contains informational data about the
// package. It is typically used in `lucy info`.
type PackageInformation struct {
	Name        string
	Brief       string
	Description string
	Author      []struct {
		Name  string
		Url   string
		Email string
	}
	Urls    []PackageUrl
	License string
}

// PackageInstallation is an optional attribution to lucytypes.Package. It is
// used for packages that are known to be installed in the local filesystem.
type PackageInstallation struct {
	Path string
}

// PackageRemote is an optional attribution to lucytypes.Package. It is used for
// packages that are known to be bound with a remote source.
type PackageRemote struct {
	Source   Source
	RemoteId string
	// The URL to download the package's specified version When package.Id.Version
	// is set to "latest" or "any", this field will be identical to LatestFileUrl.
	FileUrl string
	// The URL to get the latest version of the package.
	LatestVersionUrl string
}

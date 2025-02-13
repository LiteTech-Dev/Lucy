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

// Package is a package with all of its related information.
//
// Most of the fields other than Package.Id should be obtained from an
// external source. This includes a web API or the user's local filesystem.
//
// Here we use a combination design to extend information about a package. This is
// done by adding optional attributions to the struct.
//
// In principle, any exported function that provides a Package as return value
// should give the caller an option, or implicit option (such as the function
// naming), to specify which extra attributions are provided
type Package struct {
	Id        PackageId // Base package identifier
	Path      string
	Installed bool

	Info *PackageInformation  // Optional
	Deps *PackageDependencies // Optional
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
	Urls   []PackageUrl
	Author []struct {
		Name  string
		Url   string
		Email string
	}
	License string

	// The difference between Name and PackageId.Name is that Name is
	Name        string
	Brief       string
	Description string
}

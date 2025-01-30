package lucytypes

import "lucy/syntaxtypes"

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
	default:
		return "Unknown"
	}
}

type PackageUrl struct {
	Name string
	Type PackageUrlType
	Url  string
}

// PackageInfo is a package with all of its related information.
//
// All fields other than Id are optional.
//
// Most of the fields other than PackageInfo.Id should be obtained from an
// external source. This includes a web API or the user's local filesystem.
type PackageInfo struct {
	Id                 syntaxtypes.Package // Base package identifier
	Path               string
	Installed          bool
	Urls               []PackageUrl
	Name               string
	Description        string
	SupportedVersions  []syntaxtypes.PackageVersion
	SupportedPlatforms []syntaxtypes.Platform
}

func (p *PackageInfo) String() string {
	return p.Name
}

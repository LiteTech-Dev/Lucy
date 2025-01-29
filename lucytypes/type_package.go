package lucytypes

import "lucy/syntaxtypes"

type SourceUrlType uint8

const (
	File SourceUrlType = iota
	Homepage
	SourceCode
	Wiki
)

type PackageUrl struct {
	Name string
	Type SourceUrlType
	Url  string
}

// PackageInfo is a package with all of its related information.
//
// All fields other than Id are optional.
//
// Most of the fields other than PackageInfo.Id should be obtained from an
// external source. This includes a web API or the user's local filesystem.
type PackageInfo struct {
	Id                syntaxtypes.Package // Base package identifier
	Path              string
	Urls              []PackageUrl
	Name              string
	Description       string
	SupportedVersions []syntaxtypes.PackageVersion
}

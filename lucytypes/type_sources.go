package lucytypes

type Source uint8

const (
	CurseForge Source = iota
	Modrinth
	GitHub
	McdrSite
	UnknownSource
)

func (s Source) String() string {
	switch s {
	case CurseForge:
		return "CurseForge"
	case Modrinth:
		return "Modrinth"
	case GitHub:
		return "GitHub"
	case McdrSite:
		return "MCDR Site"
	default:
		return "Unknown Source"
	}
}

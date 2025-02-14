package lucytypes

type Source uint8

const (
	Auto Source = iota
	CurseForge
	Modrinth
	GitHub
	McdrRepo
	Unknown
)

func (s Source) String() string {
	switch s {
	case CurseForge:
		return "CurseForge"
	case Modrinth:
		return "Modrinth"
	case GitHub:
		return "GitHub"
	case McdrRepo:
		return "MCDR"
	default:
		return "Unknown Source"
	}
}

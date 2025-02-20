package lucytypes

type Source uint8

const (
	Auto Source = iota
	CurseForge
	Modrinth
	GitHub
	McdrRepo
	UnknownSource
)

func (s Source) String() string {
	switch s {
	case Auto:
		return "auto"
	case CurseForge:
		return "curseforge"
	case Modrinth:
		return "modrinth"
	case GitHub:
		return "github"
	case McdrRepo:
		return "mcdr"
	default:
		return "unknown"
	}
}

func (s Source) Title() string {
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
		return "Unknown"
	}
}

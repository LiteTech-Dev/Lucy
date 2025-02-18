package lucytypes

type SearchOptions struct {
	ShowClientPackage bool
	IndexBy           SearchIndex
}

type SearchIndex uint8

const (
	ByRelevance SearchIndex = iota
	ByDownloads
	ByNewest
)

func InputSearchIndex(input string) SearchIndex {
	switch input {
	case "relevance":
		return ByRelevance
	case "downloads":
		return ByDownloads
	case "newest":
		return ByNewest
	default:
		return ByRelevance
	}
}

func (i SearchIndex) ToModrinth() string {
	switch i {
	case ByRelevance:
		return "relevance"
	case ByDownloads:
		return "downloads"
	case ByNewest:
		return "newest"
	default:
		return "relevance"
	}
}

// func (i SearchIndex) ToCurseForge() string

type SearchResults struct {
	Source  Source
	Results []string // PackageNames
}

package lucytypes

type SearchOptions struct {
	ShowClientPackage bool
	IndexBy           SearchIndex
}

type SearchIndex string

const (
	ByRelevance = "relevance"
	ByDownloads = "downloads"
	ByNewest    = "newest"
)

func (i SearchIndex) Validate() bool {
	switch i {
	case ByRelevance, ByDownloads, ByNewest:
		return true
	default:
		return false
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

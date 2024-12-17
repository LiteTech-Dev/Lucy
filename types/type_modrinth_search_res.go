package types

import "time"

type ModrinthSearchRes struct {
	Hits []struct {
		ProjectId         string    `json:"project_id"`
		ProjectType       string    `json:"project_type"`
		Slug              string    `json:"slug"`
		Author            string    `json:"author"`
		Title             string    `json:"title"`
		Description       string    `json:"description"`
		Categories        []string  `json:"categories"`
		DisplayCategories []string  `json:"display_categories"`
		Versions          []string  `json:"versions"`
		Downloads         int       `json:"downloads"`
		Follows           int       `json:"follows"`
		IconUrl           string    `json:"icon_url"`
		DateCreated       time.Time `json:"date_created"`
		DateModified      time.Time `json:"date_modified"`
		LatestVersion     string    `json:"latest_version"`
		License           string    `json:"license"`
		ClientSide        string    `json:"client_side"`
		ServerSide        string    `json:"server_side"`
		Gallery           []string  `json:"gallery"`
		FeaturedGallery   *string   `json:"featured_gallery"`
		Color             *int      `json:"color"`
	} `json:"hits"`
	Offset    int `json:"offset"`
	Limit     int `json:"limit"`
	TotalHits int `json:"total_hits"`
}

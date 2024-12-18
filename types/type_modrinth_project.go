package types

import "time"

type ModrinthProject struct {
	ClientSide       string      `json:"client_side"`
	ServerSide       string      `json:"server_side"`
	GameVersions     []string    `json:"game_versions"`
	Id               string      `json:"id"`
	Slug             string      `json:"slug"`
	ProjectType      string      `json:"project_type"`
	Team             string      `json:"team"`
	Organization     interface{} `json:"organization"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Body             string      `json:"body"`
	BodyUrl          interface{} `json:"body_url"`
	Published        time.Time   `json:"published"`
	Updated          time.Time   `json:"updated"`
	Approved         time.Time   `json:"approved"`
	Queued           interface{} `json:"queued"`
	Status           string      `json:"status"`
	RequestedStatus  string      `json:"requested_status"`
	ModeratorMessage interface{} `json:"moderator_message"`
	License          struct {
		Id   string      `json:"id"`
		Name string      `json:"name"`
		Url  interface{} `json:"url"`
	} `json:"license"`
	Downloads            int           `json:"downloads"`
	Followers            int           `json:"followers"`
	Categories           []string      `json:"categories"`
	AdditionalCategories []interface{} `json:"additional_categories"`
	Loaders              []string      `json:"loaders"`
	Versions             []string      `json:"versions"`
	IconUrl              string        `json:"icon_url"`
	IssuesUrl            string        `json:"issues_url"`
	SourceUrl            string        `json:"source_url"`
	WikiUrl              string        `json:"wiki_url"`
	DiscordUrl           string        `json:"discord_url"`
	DonationUrls         []interface{} `json:"donation_urls"`
	Gallery              []interface{} `json:"gallery"`
	Color                int           `json:"color"`
	ThreadId             string        `json:"thread_id"`
	MonetizationStatus   string        `json:"monetization_status"`
}

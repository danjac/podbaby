package models

// SearchResult includes all search data
type SearchResult struct {
	Channels []Channel `json:"channels"`
	Podcasts []Podcast `json:"podcasts"`
}

package models

// SearchResult includes all search data
type SearchResult struct {
	Channels []Channel    `json:"channels"`
	Podcasts *PodcastList `json:"podcasts"`
}

// NewSearchResult creates a properly initialized search result
func NewSearchResult(page int) *SearchResult {
	var channels []Channel
	podcasts := NewPodcastList(page)
	return &SearchResult{channels, podcasts}
}

package models

// SearchResult includes all search data
type SearchResult struct {
	Channels []Channel    `json:"channels"`
	Podcasts *PodcastList `json:"podcasts"`
}

func NewSearchResult(page int) *SearchResult {
	channels := []Channel{}
	podcasts := NewPodcastList(page)
	return &SearchResult{channels, podcasts}
}

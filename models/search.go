package models

// SearchResult includes all search data
type SearchResult struct {
	Channels []Channel    `json:"channels"`
	Podcasts *PodcastList `json:"podcasts"`
}

func NewSearchResult() *SearchResult {
	channels := []Channel{}
	podcasts := NewPodcastList()
	return &SearchResult{channels, podcasts}
}

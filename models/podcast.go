package models

import "time"

// PodcastList is a paginated list of podcasts
type PodcastList struct {
	Podcasts []Podcast  `json:"podcasts"`
	Page     *Paginator `json:"page"`
}

// NewPodcastList creates a new PodcastList instance
func NewPodcastList(page int) *PodcastList {
	var podcasts []Podcast
	return &PodcastList{podcasts, NewPaginator(page, 0)}
}

// Podcast is a single podcast item
type Podcast struct {
	ID           int       `db:"id" json:"id"`
	ChannelID    int       `db:"channel_id" json:"channelId"`
	GUID         string    `db:"guid" json:"-"`
	Name         string    `db:"name" json:"name"`
	Image        string    `db:"image" json:"image"`
	Title        string    `db:"title" json:"title"`
	Source       string    `db:"source" json:"source"`
	Description  string    `db:"description" json:"description"`
	EnclosureURL string    `db:"enclosure_url" json:"enclosureUrl"`
	PubDate      time.Time `db:"pub_date" json:"pubDate"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
}

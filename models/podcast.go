package models

import "time"

type PodcastList struct {
	Podcasts []Podcast  `json:"podcasts"`
	Page     *Paginator `json:"page"`
}

type Podcast struct {
	ID           int       `db:"id" json:"id"`
	ChannelID    int       `db:"channel_id" json:"channelId"`
	Guid         string    `db:"guid" json:"-"`
	Name         string    `db:"name" json:"name"`
	Image        string    `db:"image" json:"image"`
	Title        string    `db:"title" json:"title"`
	Source       string    `db:"source" json:"source"`
	Description  string    `db:"description" json:"description"`
	EnclosureURL string    `db:"enclosure_url" json:"enclosureUrl"`
	PubDate      time.Time `db:"pub_date" json:"pubDate"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
}

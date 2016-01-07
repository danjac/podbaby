package models

import "time"

type PodcastList struct {
	Podcasts []Podcast `json:"podcasts"`
	Page     *Page     `json:"page"`
}

type Podcast struct {
	ID           int64     `db:"id" json:"id"`
	ChannelID    int64     `db:"channel_id" json:"channelId"`
	Name         string    `db:"name" json:"name"`
	Image        string    `db:"image" json:"image"`
	Title        string    `db:"title" json:"title"`
	Description  string    `db:"description" json:"description"`
	EnclosureURL string    `db:"enclosure_url" json:"enclosureUrl"`
	PubDate      time.Time `db:"pub_date" json:"pubDate"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
}

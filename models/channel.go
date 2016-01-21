package models

import (
	"database/sql"
)

type Channel struct {
	ID          int64          `db:"id" json:"id"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Categories  sql.NullString `db:"categories" json:"-"`
	Image       string         `db:"image" json:"image"`
	URL         string         `db:"url" json:"url"`
	Website     sql.NullString `db:"website" json:"website"`
	Podcasts    []*Podcast     `db:"-" json:"-"`
}

type ChannelDetail struct {
	Channel  *Channel  `json:"channel"`
	Podcasts []Podcast `json:"podcasts"`
	Related  []Channel `json:"relatedChannels"`
	Page     *Page     `json:"page"`
}

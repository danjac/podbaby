package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Channel struct {
	ID           int64          `db:"id" json:"id"`
	Title        string         `db:"title" json:"title"`
	Description  string         `db:"description" json:"description"`
	Image        string         `db:"image" json:"image"`
	URL          string         `db:"url" json:"url"`
	Website      sql.NullString `db:"website" json:"website"`
	CreatedAt    time.Time      `db:"created_at" json:"createdAt"`
	PubDate      pq.NullTime    `db:"pub_date" json:"pubDate"`
	Copyright    sql.NullString `db:"copyright" json:"copyright"`
	IsSubscribed bool           `db:"is_subscribed" json:"isSubscribed"`
}

type ChannelDetail struct {
	Channel  *Channel  `json:"channel"`
	Podcasts []Podcast `json:"podcasts"`
	Page     *Page     `json:"page"`
}

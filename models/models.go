package models

import "time"

type User struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type Channel struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Image       string    `db:"image" json:"image"`
	URL         string    `db:"url" json:"url"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
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
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}

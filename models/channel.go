package models

import "time"

type Channel struct {
	ID           int64     `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	Description  string    `db:"description" json:"description"`
	Image        string    `db:"image" json:"image"`
	URL          string    `db:"url" json:"url"`
	IsSubscribed bool      `db:"is_subscribed" json:"isSubscribed"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}

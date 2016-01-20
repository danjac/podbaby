package models

import "time"

type Play struct {
	PodcastID int64     `db:"podcast_id" json:"podcastId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

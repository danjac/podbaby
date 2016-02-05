package models

import "time"

// Play is a single timestamp/podcast play time
type Play struct {
	PodcastID int       `db:"podcast_id" json:"podcastId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

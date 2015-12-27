package database

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Users         UserDB
	Channels      ChannelDB
	Podcasts      PodcastDB
	Subscriptions SubscriptionDB
}

func New(db *sqlx.DB) *DB {
	return &DB{
		Users:         &defaultUserDBImpl{db},
		Channels:      &defaultChannelDBImpl{db},
		Podcasts:      &defaultPodcastDBImpl{db},
		Subscriptions: &defaultSubscriptionDBImpl{db},
	}
}

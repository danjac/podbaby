package database

import (
	"github.com/danjac/podbaby/config"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Rollback() error
	Commit() error
}

type TransactionManager interface {
	Begin() (Transaction, error)
}

type DB struct {
	*sqlx.DB
	T             TransactionManager
	Users         *UserDB
	Channels      *ChannelDB
	Categories    *CategoryDB
	Podcasts      *PodcastDB
	Bookmarks     *BookmarkDB
	Subscriptions *SubscriptionDB
	Plays         *PlayDB
}

func MustConnect(cfg *config.Config) *DB {
	return New(sqlx.MustConnect("postgres", cfg.DatabaseURL), cfg)
}

func New(db *sqlx.DB, cfg *config.Config) *DB {
	return &DB{
		DB:            db,
		Users:         newUserDB(db),
		Channels:      newChannelDB(db),
		Categories:    newCategoryDB(db),
		Podcasts:      newPodcastDB(db),
		Subscriptions: newSubscriptionDB(db),
		Bookmarks:     newBookmarkDB(db),
		Plays:         newPlayDB(db),
	}
}

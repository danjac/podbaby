package database

import (
	"github.com/danjac/podbaby/config"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Rollback() error
	Commit() error
}

type DB struct {
	*sqlx.DB
	Users         *UserDB
	Channels      *ChannelDB
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
		Podcasts:      newPodcastDB(db),
		Subscriptions: newSubscriptionDB(db),
		Bookmarks:     newBookmarkDB(db),
		Plays:         newPlayDB(db),
	}
}

func (db DB) Begin() (Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	return &DBTransaction{tx}, nil
}

type DBTransaction struct {
	*sqlx.Tx
}

func (t *DBTransaction) Commit() error {
	return t.Tx.Commit()
}

func (t *DBTransaction) Rollback() error {
	return t.Tx.Rollback()
}

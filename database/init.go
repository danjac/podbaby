package database

import (
	"github.com/danjac/podbaby/config"
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
	"path/filepath"
)

type DB struct {
	*sqlx.DB
	Users         UserDB
	Channels      ChannelDB
	Podcasts      PodcastDB
	Bookmarks     BookmarkDB
	Subscriptions SubscriptionDB
	Plays         PlayDB
}

func MustConnect(cfg *config.Config) *DB {
	db, err := New(sqlx.MustConnect("postgres", cfg.DatabaseURL), cfg)
	if err != nil {
		panic(err)
	}
	return db
}

func New(db *sqlx.DB, cfg *config.Config) (*DB, error) {

	var (
		ps  purse.Purse
		err error
	)

	if cfg.IsDev() {
		ps, err = purse.New(filepath.Join(".", "database", "queries"))
		if err != nil {
			return nil, err
		}
	} else {
		ps = Queries
	}

	return &DB{
		DB:            db,
		Users:         &defaultUserDBImpl{db, ps},
		Channels:      &defaultChannelDBImpl{db, ps},
		Podcasts:      &defaultPodcastDBImpl{db, ps},
		Subscriptions: &defaultSubscriptionDBImpl{db, ps},
		Bookmarks:     &defaultBookmarkDBImpl{db, ps},
		Plays:         &defaultPlayDBImpl{db, ps},
	}, nil
}

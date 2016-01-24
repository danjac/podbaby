package database

import (
	"github.com/danjac/podbaby/config"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

const dbContextKey = "db"

type beginner interface {
	Begin() (Transaction, error)
}

type preparer interface {
	Preparex(string) (*sqlx.Stmt, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
}

type DataHandler interface {
	sqlx.Ext
	preparer
}

type Connection interface {
	DataHandler
	beginner
}

type Transaction interface {
	DataHandler
	Rollback() error
	Commit() error
}

type DB interface {
	Connection
	Users() UserDB
	Channels() ChannelDB
	Categories() CategoryDB
	Podcasts() PodcastDB
	Bookmarks() BookmarkDB
	Subscriptions() SubscriptionDB
	Plays() PlayDB
}

type SqlDB struct {
	*sqlx.DB
}

func (db *SqlDB) Begin() (Transaction, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &SqlTransaction{tx}, nil
}

type SqlTransaction struct {
	*sqlx.Tx
}

func (db *SqlDB) Users() UserDB {
	return newUserDB()
}

func (db *SqlDB) Channels() ChannelDB {
	return newChannelDB()
}

func (db *SqlDB) Categories() CategoryDB {
	return newCategoryDB()
}

func (db *SqlDB) Bookmarks() BookmarkDB {
	return newBookmarkDB()
}

func (db *SqlDB) Podcasts() PodcastDB {
	return newPodcastDB()
}

func (db *SqlDB) Plays() PlayDB {
	return newPlayDB()
}

func (db *SqlDB) Subscriptions() SubscriptionDB {
	return newSubscriptionDB()
}

func New(cfg *config.Config) (DB, error) {
	conn, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return &SqlDB{
		DB: conn,
	}, nil
}

func NewWithContext(ctx context.Context, cfg *config.Config) (context.Context, error) {
	db, err := New(cfg)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, dbContextKey, db), nil
}

func FromContext(ctx context.Context) DB {
	return ctx.Value(dbContextKey).(DB)
}

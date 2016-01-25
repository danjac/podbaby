package store

import (
	"github.com/danjac/podbaby/config"
	"github.com/jmoiron/sqlx"
)

type closer interface {
	Close() error
}

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
	closer
}

type Transaction interface {
	DataHandler
	Rollback() error
	Commit() error
}

type Store interface {
	Conn() Connection
	Users() UserDB
	Channels() ChannelDB
	Categories() CategoryDB
	Podcasts() PodcastDB
	Bookmarks() BookmarkDB
	Subscriptions() SubscriptionDB
	Plays() PlayDB
}

type SqlStore struct {
	conn Connection
}

type SqlConnection struct {
	*sqlx.DB
}

func (store *SqlStore) Conn() Connection {
	return store.conn
}

func (conn *SqlConnection) Close() error {
	return conn.DB.Close()
}

func (conn *SqlConnection) Begin() (Transaction, error) {
	tx, err := conn.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &SqlTransaction{tx}, nil
}

type SqlTransaction struct {
	*sqlx.Tx
}

func (store *SqlStore) Users() UserDB {
	return newUserDB()
}

func (store *SqlStore) Channels() ChannelDB {
	return newChannelDB()
}

func (store *SqlStore) Categories() CategoryDB {
	return newCategoryDB()
}

func (store *SqlStore) Bookmarks() BookmarkDB {
	return newBookmarkDB()
}

func (store *SqlStore) Podcasts() PodcastDB {
	return newPodcastDB()
}

func (store *SqlStore) Plays() PlayDB {
	return newPlayDB()
}

func (store *SqlStore) Subscriptions() SubscriptionDB {
	return newSubscriptionDB()
}

func New(cfg *config.Config) (Store, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return &SqlStore{
		conn: &SqlConnection{db},
	}, nil
}

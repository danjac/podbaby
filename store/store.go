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
	Users() UserStore
	Channels() ChannelStore
	Categories() CategoryStore
	Podcasts() PodcastStore
	Bookmarks() BookmarkStore
	Subscriptions() SubscriptionStore
	Plays() PlayStore
}

type SqlStore struct {
	conn          Connection
	users         UserStore
	categories    CategoryStore
	channels      ChannelStore
	podcasts      PodcastStore
	bookmarks     BookmarkStore
	subscriptions SubscriptionStore
	plays         PlayStore
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

func (store *SqlStore) Users() UserStore {
	return store.users
}

func (store *SqlStore) Channels() ChannelStore {
	return store.channels
}

func (store *SqlStore) Categories() CategoryStore {
	return store.categories
}

func (store *SqlStore) Bookmarks() BookmarkStore {
	return store.bookmarks
}

func (store *SqlStore) Podcasts() PodcastStore {
	return store.podcasts
}

func (store *SqlStore) Plays() PlayStore {
	return store.plays
}

func (store *SqlStore) Subscriptions() SubscriptionStore {
	return newSubscriptionStore()
}

func New(cfg *config.Config) (Store, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return &SqlStore{
		conn:          &SqlConnection{db},
		categories:    newCategoryStore(),
		podcasts:      newPodcastStore(),
		bookmarks:     newBookmarkStore(),
		subscriptions: newSubscriptionStore(),
		plays:         newPlayStore(),
	}, nil
}

package store

import (
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/DATA-DOG/go-sqlmock"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
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
	Close() error
	Conn() Connection
	Users() UserStore
	Channels() ChannelStore
	Categories() CategoryStore
	Podcasts() PodcastStore
	Bookmarks() BookmarkStore
	Subscriptions() SubscriptionStore
	Plays() PlayStore
}

type sqlStore struct {
	conn          Connection
	users         UserStore
	categories    CategoryStore
	channels      ChannelStore
	podcasts      PodcastStore
	bookmarks     BookmarkStore
	subscriptions SubscriptionStore
	plays         PlayStore
}

type sqlConnection struct {
	*sqlx.DB
}

func (store *sqlStore) Close() error {
	return store.conn.Close()
}

func (store *sqlStore) Conn() Connection {
	return store.conn
}

func (conn *sqlConnection) Close() error {
	return conn.DB.Close()
}

func (conn *sqlConnection) Begin() (Transaction, error) {
	tx, err := conn.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &sqlTransaction{tx}, nil
}

type sqlTransaction struct {
	*sqlx.Tx
}

func (store *sqlStore) Users() UserStore {
	return store.users
}

func (store *sqlStore) Channels() ChannelStore {
	return store.channels
}

func (store *sqlStore) Categories() CategoryStore {
	return store.categories
}

func (store *sqlStore) Bookmarks() BookmarkStore {
	return store.bookmarks
}

func (store *sqlStore) Podcasts() PodcastStore {
	return store.podcasts
}

func (store *sqlStore) Plays() PlayStore {
	return store.plays
}

func (store *sqlStore) Subscriptions() SubscriptionStore {
	return store.subscriptions
}

func newSqlStore(db *sqlx.DB) Store {
	return &sqlStore{
		conn:          &sqlConnection{db},
		categories:    newCategoryStore(),
		channels:      newChannelStore(),
		users:         newUserStore(),
		podcasts:      newPodcastStore(),
		bookmarks:     newBookmarkStore(),
		subscriptions: newSubscriptionStore(),
		plays:         newPlayStore(),
	}
}

func New(cfg *config.Config) (Store, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	//db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(cfg.MaxDBConnections)
	return newSqlStore(db), nil
}

func NewMock() (Store, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, mock, err
	}
	dbx := sqlx.NewDb(db, "postgres")
	return newSqlStore(dbx), mock, nil
}

package api

import (
	"errors"
	"fmt"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
	"github.com/jmoiron/sqlx"
)

var errMockDBError = errors.New("Fake DB error")

type mockConnection struct {
	*sqlx.DB
}

func (conn *mockConnection) Begin() (store.Transaction, error) {
	return nil, nil
}

type mockStore struct {
	conn          store.Connection
	users         store.UserStore
	categories    store.CategoryStore
	channels      store.ChannelStore
	podcasts      store.PodcastStore
	bookmarks     store.BookmarkStore
	subscriptions store.SubscriptionStore
	plays         store.PlayStore
}

func (s *mockStore) Conn() store.Connection {
	return s.conn
}

func (s *mockStore) Users() store.UserStore {
	return s.users
}

func (s *mockStore) Channels() store.ChannelStore {
	return s.channels
}

func (s *mockStore) Categories() store.CategoryStore {
	return s.categories
}

func (s *mockStore) Bookmarks() store.BookmarkStore {
	return s.bookmarks
}

func (s *mockStore) Subscriptions() store.SubscriptionStore {
	return s.subscriptions
}

func (s *mockStore) Podcasts() store.PodcastStore {
	return s.podcasts
}

func (s *mockStore) Plays() store.PlayStore {
	return s.plays
}

type mockSelectBookmarked func(_, _ int64) (*models.PodcastList, error)

type mockPodcastStore struct {
	*store.PodcastSqlStore
	selectBookmarked mockSelectBookmarked
}

func (s *mockPodcastStore) SelectBookmarked(_ store.DataHandler, userID int64, page int64) (*models.PodcastList, error) {
	fmt.Println("Calling with userID:", userID)
	return s.selectBookmarked(userID, page)
}

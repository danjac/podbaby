package server

import (
	"fmt"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/models"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockFeedparser struct{}

func (f *mockFeedparser) Fetch(ch *models.Channel) error {
	ch.ID = 1
	ch.Title = "a new feed"
	ch.Podcasts = []*models.Podcast{
		&models.Podcast{
			ChannelID: ch.ID,
			Title:     "a random cast",
		},
	}
	return nil
}

type mockGetChannelWithNone struct {
	*database.ChannelDBReader
}

func (db *mockGetChannelWithNone) GetByURL(_ string) (*models.Channel, error) {
	return nil, &mockDBErrNoRows{}
}

type mockTransaction struct{}

func (t *mockTransaction) Commit() error {
	return nil
}

func (t *mockTransaction) Rollback() error {
	return nil
}

type mockChannelTransaction struct {
	*mockTransaction
}

func (t *mockChannelTransaction) Create(_ *models.Channel) error {
	return nil
}

func (t *mockChannelTransaction) AddCategories(_ *models.Channel) error {
	return nil
}

func (t *mockChannelTransaction) AddPodcasts(_ *models.Channel) error {
	return nil
}

type mockChannelWriter struct{}

func (w *mockChannelWriter) Begin() (database.ChannelTransaction, error) {
	return &mockChannelTransaction{}, nil
}

type mockSubscriptionWriter struct {
	*database.SubscriptionDBWriter
}

func (w *mockSubscriptionWriter) Create(_, _ int64) error {
	return nil
}

func TestAddChannelIfNew(t *testing.T) {

	user := &models.User{
		ID: 1,
	}

	getContext = mockGetContextWithUser(user)
	url := "http://joeroganexp.joerogan.libsynpro.com/rss"

	r := strings.NewReader(fmt.Sprintf(`{ "url": "%s" }`, url))
	req, _ := http.NewRequest("GET", "/", r)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Channels: &database.ChannelDB{
				ChannelReader: &mockGetChannelWithNone{},
				ChannelWriter: &mockChannelWriter{},
			},
			Subscriptions: &database.SubscriptionDB{
				SubscriptionWriter: &mockSubscriptionWriter{},
			},
		},
		Feedparser: &mockFeedparser{},
		Render:     render.New(),
	}

	if err := addChannel(s, w, req); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusCreated {
		t.Fatal("Should be a new channel")
	}
	body := w.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`"url":"%s"`, url)) {
		t.Fatal("Channel should be returned")
	}

}

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

type mockCreateChannel struct {
	*database.ChannelDBWriter
	isCalled bool
}

func (db *mockCreateChannel) Create(_ *models.Channel) error {
	db.isCalled = true
	return nil
}

type mockGetChannelWithNone struct {
	*database.ChannelDBReader
}

func (db *mockGetChannelWithNone) GetByURL(_ string) (*models.Channel, error) {
	return nil, &mockDBErrNoRows{}
}

type mockCreateSubscription struct {
	*database.SubscriptionDBWriter
	isCalled bool
}

func (db *mockCreateSubscription) Create(_, _ int64) error {
	db.isCalled = true
	return nil
}

type mockCreatePodcast struct {
	*database.PodcastDBWriter
	isCalled bool
}

func (db *mockCreatePodcast) Create(_ *models.Podcast) error {
	db.isCalled = true
	return nil
}

type mockTransaction struct{}

func (t *mockTransaction) Commit() error {
	return nil
}

func (t *mockTransaction) Rollback() error {
	return nil
}

type mockTransactionManager struct{}

func (tm *mockTransactionManager) Begin() (database.Transaction, error) {
	return &mockTransaction{}, nil
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
			T: &mockTransactionManager{},
			Podcasts: &database.PodcastDB{
				PodcastWriter: &mockCreatePodcast{},
			},
			Channels: &database.ChannelDB{
				ChannelReader: &mockGetChannelWithNone{},
				ChannelWriter: &mockCreateChannel{},
			},
			Subscriptions: &database.SubscriptionDB{
				SubscriptionWriter: &mockCreateSubscription{},
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

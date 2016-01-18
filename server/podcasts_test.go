package server

import (
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/models"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockGetLatestPodcasts struct {
	*database.PodcastDBReader
}

func (db *mockGetLatestPodcasts) SelectSubscribed(_, _ int64) (*models.PodcastList, error) {
	return &models.PodcastList{
		Podcasts: []models.Podcast{
			models.Podcast{
				ID:    1,
				Title: "My subscribed podcast",
			},
		},
	}, nil
}

func (db *mockGetLatestPodcasts) SelectAll(_ int64) (*models.PodcastList, error) {
	return &models.PodcastList{
		Podcasts: []models.Podcast{
			models.Podcast{
				ID:    1,
				Title: "My subscribed podcast",
			},
			models.Podcast{
				ID:    2,
				Title: "Another podcast",
			},
		},
	}, nil
}

func TestLatestPodcastsIfLoggedIn(t *testing.T) {

	user := &models.User{
		ID: 1,
	}

	getContext = mockGetContextWithUser(user)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Podcasts: &database.PodcastDB{
				PodcastReader: &mockGetLatestPodcasts{},
			},
		},
		Render: render.New(),
	}

	if err := getLatestPodcasts(s, w, req); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}
	content := w.Body.String()

	if strings.Count(content, "\"id\"") != 1 {
		t.Fatal("Should contain only own subscriptions: ", content)
	}

}

func TestLatestPodcastsIfNotLoggedIn(t *testing.T) {

	getContext = mockGetEmptyContext()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Podcasts: &database.PodcastDB{
				PodcastReader: &mockGetLatestPodcasts{},
			},
		},
		Render: render.New(),
	}

	if err := getLatestPodcasts(s, w, req); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}
	content := w.Body.String()

	if strings.Count(content, "\"id\"") != 2 {
		t.Fatal("Should contain all podcasts: ", content)
	}

}

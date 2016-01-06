package server

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/models"
	"github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"testing"
)

var errMockDBError = errors.New("Fake DB error")

type mockPodcastsDB struct {
	hasError bool
}

func (db *mockPodcastsDB) SelectSubscribed(_, _ int64) (*models.PodcastList, error) {
	return nil, nil
}

func (db *mockPodcastsDB) SelectPlayed(_, _ int64) (*models.PodcastList, error) {
	return nil, nil
}

func (db *mockPodcastsDB) SelectByChannelID(_, _, _ int64) (*models.PodcastList, error) {
	return nil, nil
}

func (db *mockPodcastsDB) Search(_ string, _ int64) ([]models.Podcast, error) {
	return nil, nil
}

func (db *mockPodcastsDB) SearchBookmarked(_ string, _ int64) ([]models.Podcast, error) {
	return nil, nil
}

func (db *mockPodcastsDB) SearchByChannelID(_ string, _, _ int64) ([]models.Podcast, error) {
	return nil, nil
}

func (db *mockPodcastsDB) Create(_ *models.Podcast) error { return nil }

func (db *mockPodcastsDB) SelectBookmarked(userID, page int64) (*models.PodcastList, error) {
	if db.hasError {
		return nil, errMockDBError
	}
	result := &models.PodcastList{}
	result.Podcasts = []models.Podcast{
		models.Podcast{
			ID:    100,
			Title: "testing",
		},
	}
	result.Page = &models.Page{}
	return result, nil
}

func TestGetBookmarksIfNotOk(t *testing.T) {

	user := &models.User{
		ID: 10,
	}

	getContext = mockGetContextWithUser(user)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Podcasts: &mockPodcastsDB{hasError: true},
		},
		Log:    logrus.New(),
		Render: render.New(),
	}
	s.getBookmarks(w, req)

	if w.Code == http.StatusOK {
		t.Fail()
	}
}

func TestGetBookmarksIfOk(t *testing.T) {

	user := &models.User{
		ID: 10,
	}

	getContext = mockGetContextWithUser(user)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Podcasts: &mockPodcastsDB{},
		},
		Render: render.New(),
	}
	s.getBookmarks(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

}

package server

import (
	"errors"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/models"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var errMockDBError = errors.New("Fake DB error")

type mockFailPodcastReader struct {
	*database.PodcastDBReader
}

func (db *mockFailPodcastReader) SelectBookmarked(_, _ int64) (*models.PodcastList, error) {
	return nil, errMockDBError
}

type mockOkPodcastReader struct {
	*database.PodcastDBReader
}

func (db *mockOkPodcastReader) SelectBookmarked(_, _ int64) (*models.PodcastList, error) {
	return &models.PodcastList{}, nil
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
			Podcasts: &database.PodcastDB{
				PodcastReader: &mockFailPodcastReader{},
			},
		},
		Render: render.New(),
	}
	if err := getBookmarks(s, w, req); err == nil {
		t.Fatal("This should return an error")
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
			Podcasts: &database.PodcastDB{
				PodcastReader: &mockOkPodcastReader{},
			},
		},
		Render: render.New(),
	}
	if err := getBookmarks(s, w, req); err != nil {
		t.Fatal("Should not return an error")
	}

	if w.Code != http.StatusOK {
		t.Fail()
	}

}

package api

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLatestPodcastsIfLoggedIn(t *testing.T) {

	user := &models.User{
		ID: 1,
	}

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()

	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	c.Set(authenticatorContextKey, &fakeAuthenticator{user})
	c.Set(storeContextKey, s)

	rows := sqlmock.NewRows([]string{""}).AddRow(1)
	mock.ExpectQuery(`^SELECT SUM*`).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{
		"id", "title", "enclosure_url", "description",
		"channel_id", "title", "image", "pub_date", "source",
	}).AddRow(1, "test", "test,mp3", "test", 2, "testing", "test.jpg", time.Now(), "")
	mock.ExpectQuery("^SELECT p.id, (.+) FROM podcasts p*").WillReturnRows(rows)

	if err := getLatestPodcasts(c); err != nil {
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

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()

	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	c.Set(authenticatorContextKey, &fakeAuthenticator{nil})
	c.Set(storeContextKey, s)

	rows := sqlmock.NewRows([]string{""}).AddRow(1)
	mock.ExpectQuery(`^SELECT reltuples*`).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{
		"id", "title", "enclosure_url", "description",
		"channel_id", "title", "image", "pub_date", "source",
	}).AddRow(1, "test", "test,mp3", "test", 2, "testing", "test.jpg", time.Now(), "")
	mock.ExpectQuery("^SELECT p.id, (.+) FROM podcasts p*").WillReturnRows(rows)

	if err := getLatestPodcasts(c); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}
	content := w.Body.String()

	if strings.Count(content, "\"id\"") != 1 {
		t.Fatal("Should contain all podcasts", content)
	}

}

package api

import (
	"fmt"
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

func TestGetBookmarksIfOk(t *testing.T) {

	user := &models.User{
		ID: 10,
	}

	req, _ := http.NewRequest("GET", "/5/", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	/*
		r := e.Router()
		r.Add(echo.GET, "/:id/", nil, e)
		r.Find(echo.GET, "/5/", c)
	*/

	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{""}).AddRow(1)
	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM bookmarks*`).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{
		"id", "title", "enclosure_url", "description",
		"channel_id", "title", "image", "pub_date", "source",
	}).AddRow(1, "test", "test,mp3", "test", 2, "testing", "test.jpg", time.Now(), "")
	mock.ExpectQuery("^SELECT p.id, (.+) FROM podcasts p*").WillReturnRows(rows)

	c.Set(userContextKey, user)
	c.Set(storeContextKey, s)

	if err := getBookmarks(c); err != nil {
		t.Fatal(err)
		t.Fatal("This should not return an error")
	}

	body := w.Body.String()
	if !strings.Contains(body, `"title":"testing"`) {
		t.Fatal("Should contain title 'testing'")
	}

}

func TestGetBookmarksIfNotOk(t *testing.T) {

	user := &models.User{
		ID: 10,
	}

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	e := echo.New()

	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}
	result := fmt.Errorf("some error")
	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM bookmarks*`).WillReturnError(result)

	c.Set(userContextKey, user)
	c.Set(storeContextKey, s)

	if err := getBookmarks(c); err == nil {
		t.Fatal("This should not return an error")
	}

}

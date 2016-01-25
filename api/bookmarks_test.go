package api

import (
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBookmarksIfNotOk(t *testing.T) {

	user := &models.User{
		ID: 10,
	}

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	e := echo.New()

	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	var mockSelectBookmarked = func(_, _ int64) (*models.PodcastList, error) {
		log.Println("CALLIng MOCK SELECT BOOKMARKS")
		return nil, errMockDBError
	}

	s := &mockStore{
		conn: &mockConnection{},
		podcasts: &mockPodcastStore{
			selectBookmarked: mockSelectBookmarked,
		},
	}

	c.Set(userContextKey, user)
	c.Set(storeContextKey, s)

	if err := getBookmarks(c); err == nil {
		t.Fatal("This should return an error")
	}

}

/*

func TestGetBookmarksIfOk(t *testing.T) {

	user := &models.User{
		ID: 10,
	}

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	e := echo.New()

	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	var getBookmarks = func(_ int64) (*models.PodcastList, error) {
		return &models.PodcastList{}, nil
	}

	s := &mockStore{
		bookmarks: mockBookmarkStore{
			getBookmarks: getBookmarks,
		},
	}

	c.Set(userContextKey, user)
	c.Set(storeContextKey, s)

	if err := getBookmarks(c); err == nil {
		t.Fatal("This should return an error")
	}

}
*/

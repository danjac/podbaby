package api

import (
	"fmt"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/DATA-DOG/go-sqlmock"
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

func TestAddChannelIfNew(t *testing.T) {

	user := &models.User{
		ID: 1,
	}

	url := "http://joeroganexp.joerogan.libsynpro.com/rss"

	r := strings.NewReader(fmt.Sprintf(`{ "url": "%s" }`, url))
	req, _ := http.NewRequest("GET", "/", r)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	empty := sqlmock.NewRows([]string{})
	newChannel := sqlmock.NewRows([]string{""}).AddRow(1)
	newPodcast := sqlmock.NewRows([]string{""}).AddRow(1)
	newSub := sqlmock.NewResult(1, 1)

	mock.ExpectQuery("^SELECT (.+) FROM channels*").WillReturnRows(empty)
	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT upsert_channel (.+)*").WillReturnRows(newChannel)
	mock.ExpectPrepare("^SELECT insert_podcast*").WillReturnError(nil)
	mock.ExpectQuery("^SELECT insert_podcast*").WillReturnRows(newPodcast)
	mock.ExpectExec("^INSERT INTO subscriptions*").WillReturnResult(newSub)
	mock.ExpectCommit()

	c.Set(storeContextKey, s)
	c.Set(userContextKey, user)
	c.Set(feedparserContextKey, &mockFeedparser{})

	c.Request().Header.Set(echo.ContentType, "application/json")

	if err := addChannel(c); err != nil {
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

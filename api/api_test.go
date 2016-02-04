package api

import (
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/cache"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeAuthenticator struct {
	user *models.User
}

func (a *fakeAuthenticator) authenticate(c *echo.Context) (*models.User, error) {
	return a.user, nil
}

type fakeCache struct{}

func (c *fakeCache) Delete(string) error { return nil }

func (c *fakeCache) Get(_ string, _ time.Duration, _ interface{}, fn cache.Setter) error {
	return fn()
}

type fakeSession struct{}

func (s *fakeSession) Read(_ *echo.Context, _ string) (interface{}, error) {
	return nil, nil
}

func (s *fakeSession) Write(_ *echo.Context, _ string, _ interface{}) error {
	return nil
}

type fakeEmptySession struct {
	*fakeSession
}

func (s *fakeEmptySession) Read(_ *echo.Context, _ string) (interface{}, error) {
	return 0, nil
}

type fakeNonEmptySession struct {
	*fakeSession
	value interface{}
}

func (s *fakeNonEmptySession) Read(_ *echo.Context, _ string) (interface{}, error) {
	return s.value, nil
}

func TestDefaultAuthenticatorIfIsUser(t *testing.T) {

	req, _ := http.NewRequest("GET", "/5/", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	userID := 1

	session := &fakeNonEmptySession{value: userID}

	c.Set(sessionContextKey, session)

	a := &defaultAuthenticator{}

	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at"}).AddRow(1, "tester", "tester@gmail.com", "testpass", time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM users*").WillReturnRows(rows)

	c.Set(storeContextKey, s)

	user, err := a.authenticate(c)
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("We should return a valid user")
	}

	if c.Get(userContextKey) == nil {
		t.Fatal("User should be added to context")
	}

}
func TestDefaultAuthenticatorIfEmpty(t *testing.T) {

	req, _ := http.NewRequest("GET", "/5/", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	c.Set(sessionContextKey, &fakeEmptySession{})

	a := &defaultAuthenticator{}

	user, err := a.authenticate(c)
	if err != nil {
		t.Fatal(err)
	}

	if user != nil {
		t.Fatal("We should have a nil user if session empty")
	}

}

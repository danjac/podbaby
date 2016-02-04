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

type fakeSession struct {
	value int
}

func (s *fakeSession) write(_ *echo.Context, _ string, _ interface{}) error { return nil }
func (s *fakeSession) read(_ *echo.Context, _ string, _ interface{}) error  { return nil }
func (s *fakeSession) readInt(_ *echo.Context, _ string) (int, error)       { return s.value, nil }

func TestDefaultAuthenticatorMemoized(t *testing.T) {

	// should just return user if already in context

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	existingUser := &models.User{ID: 1}
	c.Set(userContextKey, existingUser)

	a := &defaultAuthenticator{}
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

	if getUser(c).ID != existingUser.ID {
		t.Fatal("Should be the same user")
	}
}

func TestDefaultAuthenticatorIfIsUser(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	session := &fakeSession{value: 1}

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

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	session := &fakeSession{value: 0}
	c.Set(sessionContextKey, session)

	a := &defaultAuthenticator{}

	user, err := a.authenticate(c)
	if err != nil {
		t.Fatal(err)
	}

	if user != nil {
		t.Fatal("We should have a nil user if session empty")
	}

}

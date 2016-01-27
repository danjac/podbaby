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
)

func TestIsEmailIfLoggedInAndAnotherEmailExists(t *testing.T) {

	user := &models.User{
		ID: 2,
	}
	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{""}).AddRow(1)

	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM users*`).WillReturnRows(rows)
	req, _ := http.NewRequest("GET", "/?email=test@gmail.com", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	c.Set(storeContextKey, s)
	c.Set(userContextKey, user)
	c.Set(authenticatorContextKey, &fakeAuthenticator{user})

	if err := isEmail(c); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}

	if !strings.Contains(w.Body.String(), "true") {
		t.Fatalf("should return true if someone else's email:%s", w.Body.String())
	}
}

func TestIsEmailIfLoggedInAndOwnEmailExists(t *testing.T) {

	user := &models.User{
		ID: 2,
	}
	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{""}).AddRow(0)

	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM users*`).WillReturnRows(rows)
	req, _ := http.NewRequest("GET", "/?email=test@gmail.com", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	c.Set(storeContextKey, s)
	c.Set(userContextKey, user)
	c.Set(authenticatorContextKey, &fakeAuthenticator{user})

	if err := isEmail(c); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}

	if !strings.Contains(w.Body.String(), "false") {
		t.Fatalf("should return false if user's own email:%s", w.Body.String())
	}
}

func TestIsEmailIfNotLoggedInAndEmailExists(t *testing.T) {

	user := &models.User{
		ID: 2,
	}
	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{""}).AddRow(1)

	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM users*`).WillReturnRows(rows)
	req, _ := http.NewRequest("GET", "/?email=test@gmail.com", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	c.Set(storeContextKey, s)
	c.Set(userContextKey, user)
	c.Set(authenticatorContextKey, &fakeAuthenticator{nil})

	if err := isEmail(c); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}

	if !strings.Contains(w.Body.String(), "true") {
		t.Fatalf("should return true if email taken %s", w.Body.String())
	}
}

func TestIsEmailIfNotLoggedInAndEmailDoesNotExist(t *testing.T) {

	user := &models.User{
		ID: 2,
	}
	s, mock, err := store.NewMock()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{""}).AddRow(0)

	mock.ExpectQuery(`^SELECT COUNT\(id\) FROM users*`).WillReturnRows(rows)
	req, _ := http.NewRequest("GET", "/?email=test@gmail.com", nil)
	w := httptest.NewRecorder()

	e := echo.New()
	c := echo.NewContext(req, echo.NewResponse(w, e), e)

	c.Set(storeContextKey, s)
	c.Set(userContextKey, user)
	c.Set(authenticatorContextKey, &fakeAuthenticator{nil})

	if err := isEmail(c); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}

	if !strings.Contains(w.Body.String(), "false") {
		t.Fatalf("should return false if email not taken %s", w.Body.String())
	}
}

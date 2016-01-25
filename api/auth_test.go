package api

/*
import (
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/models"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockIsEmailExistsUserReader struct {
	*database.UserDBReader
}

func (db *mockIsEmailExistsUserReader) IsEmail(email string, userID int64) (bool, error) {
	if email == "test@gmail.com" {
		return userID != 1, nil
	}
	return false, nil
}

func TestIsEmailIfLoggedInAndAnotherEmailExists(t *testing.T) {

	user := &models.User{
		ID: 2,
	}

	getContext = mockGetContextWithUser(user)

	req, _ := http.NewRequest("GET", "/?email=test@gmail.com", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Users: &database.UserDB{
				UserReader: &mockIsEmailExistsUserReader{},
			},
		},
		Render: render.New(),
	}

	if err := isEmail(s, w, req); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}

	if !strings.Contains(w.Body.String(), "true") {
		t.Fatalf("should return false if user's own email:%s", w.Body.String())
	}
}

func TestIsEmailIfLoggedInAndEmailExists(t *testing.T) {

	user := &models.User{
		ID: 1,
	}

	getContext = mockGetContextWithUser(user)

	req, _ := http.NewRequest("GET", "/?email=test@gmail.com", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Users: &database.UserDB{
				UserReader: &mockIsEmailExistsUserReader{},
			},
		},
		Render: render.New(),
	}

	if err := isEmail(s, w, req); err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}

	if !strings.Contains(w.Body.String(), "false") {
		t.Fatalf("should return false if user's own email:%s", w.Body.String())
	}
}

func TestIsEmailIfNotLoggedInAndEmailNotExists(t *testing.T) {

	getContext = mockGetContext(make(map[string]interface{}))

	req, _ := http.NewRequest("GET", "/?email=test2@gmail.com", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Users: &database.UserDB{
				UserReader: &mockIsEmailExistsUserReader{},
			},
		},
		Render: render.New(),
	}

	if err := isEmail(s, w, req); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}
	if !strings.Contains(w.Body.String(), "false") {
		t.Fatal("should return false if used email")
	}

}
func TestIsEmailIfNotLoggedInAndEmailExists(t *testing.T) {

	getContext = mockGetContext(make(map[string]interface{}))

	req, _ := http.NewRequest("GET", "/?email=test@gmail.com", nil)
	w := httptest.NewRecorder()
	s := &Server{
		DB: &database.DB{
			Users: &database.UserDB{
				UserReader: &mockIsEmailExistsUserReader{},
			},
		},
		Render: render.New(),
	}

	if err := isEmail(s, w, req); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK {
		t.Fatal("Should return a 200 OK")
	}
	if !strings.Contains(w.Body.String(), "true") {
		t.Fatal("should return true if used email")
	}

}
*/

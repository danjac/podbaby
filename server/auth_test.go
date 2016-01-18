package server

import (
	"github.com/danjac/podbaby/database"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockIsEmailExistsUserReader struct {
	*database.UserDBReader
}

func (db *mockIsEmailExistsUserReader) IsEmail(email string, userID int64) (bool, error) {
	return true, nil
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

}

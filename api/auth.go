package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/danjac/podbaby/models"
	"github.com/gorilla/context"
)

const (
	cookieUserID  = "userid"
	userKey       = "user"
	cookieTimeout = 24
)

// authentication methods

func getUser(r *http.Request) (*models.User, bool) {
	val, ok := context.GetOk(r, userKey)
	if !ok {
		return nil, false
	}
	return val.(*models.User), true
}

func (s *Server) setAuthCookie(w http.ResponseWriter, userID int64) {

	if encoded, err := s.Cookie.Encode(cookieUserID, userID); err == nil {
		cookie := &http.Cookie{
			Name:    cookieUserID,
			Value:   encoded,
			Expires: time.Now().Add(time.Hour * cookieTimeout),
			//Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}
}

func (s *Server) requireAuth(fn http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if user already set elsewhere
		if _, ok := getUser(r); ok {
			fn(w, r)
			return
		}
		// get user from cookie
		user, err := s.getUserFromCookie(r)
		if err != nil {
			s.abort(w, r, err)
			return
		}
		// all ok...
		context.Set(r, userKey, user)
		fn(w, r)
	})

}

func (s *Server) getUserFromCookie(r *http.Request) (*models.User, error) {

	cookie, err := r.Cookie(cookieUserID)
	if err != nil {
		return nil, HTTPError{http.StatusUnauthorized, err}
	}

	var userID int64

	if err := s.Cookie.Decode(cookieUserID, cookie.Value, &userID); err != nil {
		return nil, HTTPError{http.StatusUnauthorized, err}
	}

	if userID == 0 {
		return nil, HTTPError{http.StatusUnauthorized, errors.New("Cookie is empty")}
	}

	user, err := s.DB.Users.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, HTTPError{http.StatusUnauthorized, errors.New("No user found for this ID")}
		}
		return nil, err
	}
	return user, nil

}

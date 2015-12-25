package server

import (
	"database/sql"
	"errors"
	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/models"
	"net/http"
	"strconv"
	"time"
)

const cookieUserID = "userid"

// auth routes

func setAuthCookie(w http.ResponseWriter, userID int64) {
	cookie := &http.Cookie{
		Name:    cookieUserID,
		Value:   strconv.FormatInt(userID, 10),
		Expires: time.Now().Add(time.Hour),
		//Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func (s *Server) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	// log in here, set cookie, return username
	cookie, err := r.Cookie(cookieUserID)

	if err != nil {
		s.Abort(w, r, HTTPError{http.StatusUnauthorized, err})
		return
	}

	if cookie.Value == "" || cookie.Value == "0" {
		s.Abort(w, r, HTTPError{http.StatusUnauthorized, errors.New("Unauthenticated")})
		return
	}

	user, err := s.DB.Users.GetByID(cookie.Value)
	if err != nil {
		s.Abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, user)

}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.Signup{}

	if err := decoders.Decode(r, decoder); r != nil {
		s.Abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	if exists, _ := s.DB.Users.IsEmail(decoder.Email); exists {
		s.Abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Email taken")})
		return
	}

	if exists, _ := s.DB.Users.IsName(decoder.Name); exists {
		s.Abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Name taken")})
		return
	}

	// make new user

	user := &models.User{
		Name:  decoder.Name,
		Email: decoder.Email,
	}

	if err := user.SetPassword(decoder.Password); err != nil {
		s.Abort(w, r, err)
		return
	}

	if err := s.DB.Users.Create(user); err != nil {
		s.Abort(w, r, err)
		return
	}
	setAuthCookie(w, user.ID)
	// tbd: no need to return user!
	s.Render.JSON(w, http.StatusCreated, user)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.Login{}

	if err := decoders.Decode(r, decoder); err != nil {
		s.Abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {

		if err == sql.ErrNoRows {
			s.Abort(w, r, HTTPError{http.StatusBadRequest, errors.New("no user found")})
			return
		}
		s.Abort(w, r, err)
		return
	}

	if !user.CheckPassword(decoder.Password) {
		s.Abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Invalid password")})
		return
	}
	// login user
	setAuthCookie(w, user.ID)

	// tbd: no need to return user!
	s.Render.JSON(w, http.StatusOK, user)

}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	setAuthCookie(w, 0)
	s.Render.Text(w, http.StatusOK, "Logged out")
}

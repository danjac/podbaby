package api

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomPassword(length int) string {
	b := make([]byte, length)
	numChars := len(passwordChars)
	for i := range b {
		b[i] = passwordChars[rand.Intn(numChars)]
	}
	return string(b)
}

func (s *Server) recoverPassword(w http.ResponseWriter, r *http.Request) {
	// generate a temp password
	decoder := &decoders.RecoverPassword{}

	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {

		if err == sql.ErrNoRows {
			s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("no user found")})
			return
		}
		s.abort(w, r, err)
		return
	}

	tempPassword := generateRandomPassword(6)

	user.SetPassword(tempPassword)

	if err := s.DB.Users.UpdatePassword(user.Password, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}

	// send email to user
	// TBD: add email config to server config
	go func(r *http.Request, user *models.User, tempPassword string) {

		msg := fmt.Sprintf(`Hi %s,
We've reset your password so you can sign back in again!

Here is your new temporary password:

%s

You can login here:

%s/#/login/

Change your password as soon as possible!

Thanks,

PodBaby
    `, user.Name, tempPassword, r.Host)

		s.Log.Info(msg)

		err := smtp.SendMail(
			"mail.localhost:25",
			nil, // auth
			"sender@podbaby.me",
			[]string{user.Email},
			[]byte(msg),
		)

		if err != nil {
			s.Log.Error(err)
		}

	}(r, user, tempPassword)

	s.Render.Text(w, http.StatusOK, "password sent")

}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.Signup{}

	if err := decoders.Decode(r, decoder); r != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	if exists, _ := s.DB.Users.IsEmail(decoder.Email, 0); exists {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Email taken")})
		return
	}

	if exists, _ := s.DB.Users.IsName(decoder.Name); exists {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Name taken")})
		return
	}

	// make new user

	user := &models.User{
		Name:  decoder.Name,
		Email: decoder.Email,
	}

	if err := user.SetPassword(decoder.Password); err != nil {
		s.abort(w, r, err)
		return
	}

	if err := s.DB.Users.Create(user); err != nil {
		s.abort(w, r, err)
		return
	}
	s.setAuthCookie(w, user.ID)
	// tbd: no need to return user!
	s.Render.JSON(w, http.StatusCreated, user)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	decoder := &decoders.Login{}
	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {

		if err == sql.ErrNoRows {
			s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("no user found")})
			return
		}
		s.abort(w, r, err)
		return
	}

	if !user.CheckPassword(decoder.Password) {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Invalid password")})
		return
	}
	// login user
	s.setAuthCookie(w, user.ID)

	// tbd: no need to return user!
	s.Render.JSON(w, http.StatusOK, user)

}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	s.setAuthCookie(w, 0)
	s.Render.Text(w, http.StatusOK, "Logged out")
}

func (s *Server) changeEmail(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	decoder := &decoders.NewEmail{}
	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}
	// does this email exist?
	if exists, _ := s.DB.Users.IsEmail(decoder.Email, user.ID); exists {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Email taken")})
		return
	}

	if err := s.DB.Users.UpdateEmail(decoder.Email, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "email updated")
}

func (s *Server) changePassword(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	decoder := &decoders.NewPassword{}
	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	// validate old password first

	if !user.CheckPassword(decoder.OldPassword) {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Invalid password")})
		return
	}
	user.SetPassword(decoder.NewPassword)

	if err := s.DB.Users.UpdatePassword(user.Password, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "password updated")
}

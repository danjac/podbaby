package server

import (
	"database/sql"
	"math/rand"
	"net/http"
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
		s.abort(w, r, err)
		return
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {

		if err == sql.ErrNoRows {
			errors := decoders.Errors{
				"identifier": "No user found matching this email or name",
			}
			s.abort(w, r, errors)
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

	data := map[string]string{
		"name":         user.Name,
		"tempPassword": tempPassword,
		"host":         r.Host,
	}

	go func(to string, data map[string]string) {

		err := s.Mailer.SendFromTemplate(
			"services@podbaby.me",
			[]string{to},
			"Your new password",
			"recover_password.tmpl",
			data,
		)
		if err != nil {
			s.Log.Error(err)
		}

	}(user.Email, data)

	s.Render.Text(w, http.StatusOK, "password sent")

}

func (s *Server) isName(w http.ResponseWriter, r *http.Request) {

	var (
		err    error
		exists bool
	)

	name := r.FormValue("name")
	if name == "" {
		s.Render.Text(w, http.StatusBadRequest, "name param required")
		return
	}

	if exists, err = s.DB.Users.IsName(name); err != nil {
		s.abort(w, r, err)
		return
	}

	s.Render.JSON(w, http.StatusOK, exists)
}

func (s *Server) isEmail(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		exists bool
		userID int64
	)

	email := r.FormValue("email")
	if email == "" {
		s.Render.Text(w, http.StatusBadRequest, "email param required")
		return
	}

	if user, err := s.getUserFromCookie(r); err == nil {
		userID = user.ID
	}

	if exists, err = s.DB.Users.IsEmail(email, userID); err != nil {
		s.abort(w, r, err)
		return
	}

	s.Render.JSON(w, http.StatusOK, exists)

}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.Signup{}

	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, err)
		return
	}

	errors := make(decoders.Errors)
	if exists, _ := s.DB.Users.IsEmail(decoder.Email, 0); exists {
		errors["email"] = "This email address is taken"
	}

	if exists, _ := s.DB.Users.IsName(decoder.Name); exists {
		errors["name"] = "This name is taken"
	}

	if len(errors) > 0 {
		s.abort(w, r, errors)
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
		s.abort(w, r, err)
		return
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {

		if err == sql.ErrNoRows {
			errors := decoders.Errors{
				"identifier": "No user found matching this name or email",
			}
			s.abort(w, r, errors)
			return
		}
		s.abort(w, r, err)
		return
	}

	if !user.CheckPassword(decoder.Password) {
		errors := decoders.Errors{
			"password": "Your password is invalid",
		}
		s.abort(w, r, errors)
		return
	}

	// get bookmarks & subscriptions
	if user.Bookmarks, err = s.DB.Bookmarks.SelectByUserID(user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	if user.Subscriptions, err = s.DB.Subscriptions.SelectByUserID(user.ID); err != nil {
		s.abort(w, r, err)
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
		s.abort(w, r, err)
		return
	}

	// does this email exist?
	if exists, _ := s.DB.Users.IsEmail(decoder.Email, user.ID); exists {
		errors := decoders.Errors{
			"email": "This email address is taken",
		}
		s.abort(w, r, errors)
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
		s.abort(w, r, err)
		return
	}

	// validate old password first

	if !user.CheckPassword(decoder.OldPassword) {
		errors := decoders.Errors{
			"oldPassword": "Your old password was not recognized",
		}
		s.abort(w, r, errors)
		return
	}
	user.SetPassword(decoder.NewPassword)

	if err := s.DB.Users.UpdatePassword(user.Password, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "password updated")
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	if err := s.DB.Users.DeleteUser(user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	// log user out
	s.setAuthCookie(w, 0)
	s.Render.Text(w, http.StatusOK, "account deleted")
}

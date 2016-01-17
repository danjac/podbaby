package server

import (
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

func recoverPassword(s *Server, w http.ResponseWriter, r *http.Request) error {
	// generate a temp password
	decoder := &decoders.RecoverPassword{}

	if err := decoders.Decode(r, decoder); err != nil {
		return err
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {

		if isErrNoRows(err) {
			errors := decoders.Errors{
				"identifier": "No user found matching this email or name",
			}
			return errors
		}
		return err
	}

	tempPassword := generateRandomPassword(6)

	user.SetPassword(tempPassword)

	if err := s.DB.Users.UpdatePassword(user.Password, user.ID); err != nil {
		return err
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

	return s.Render.Text(w, http.StatusOK, "password sent")

}

func isName(s *Server, w http.ResponseWriter, r *http.Request) error {

	var (
		err    error
		exists bool
	)

	name := r.FormValue("name")
	if name == "" {
		return s.Render.Text(w, http.StatusBadRequest, "name param required")
	}

	if exists, err = s.DB.Users.IsName(name); err != nil {
		return err
	}

	return s.Render.JSON(w, http.StatusOK, exists)
}

func isEmail(s *Server, w http.ResponseWriter, r *http.Request) error {
	var (
		err    error
		exists bool
		userID int64
	)

	email := r.FormValue("email")
	if email == "" {
		return s.Render.Text(w, http.StatusBadRequest, "email param required")
	}

	if user, err := s.getUserFromCookie(r); err == nil {
		userID = user.ID
	}

	if exists, err = s.DB.Users.IsEmail(email, userID); err != nil {
		return err
	}

	return s.Render.JSON(w, http.StatusOK, exists)

}

func signup(s *Server, w http.ResponseWriter, r *http.Request) error {

	decoder := &decoders.Signup{}

	if err := decoders.Decode(r, decoder); err != nil {
		return err
	}

	errors := make(decoders.Errors)
	if exists, _ := s.DB.Users.IsEmail(decoder.Email, 0); exists {
		errors["email"] = "This email address is taken"
	}

	if exists, _ := s.DB.Users.IsName(decoder.Name); exists {
		errors["name"] = "This name is taken"
	}

	if len(errors) > 0 {
		return errors
	}

	// make new user

	user := &models.User{
		Name:  decoder.Name,
		Email: decoder.Email,
	}

	if err := user.SetPassword(decoder.Password); err != nil {
		return err
	}

	if err := s.DB.Users.Create(user); err != nil {
		return err
	}
	s.setAuthCookie(w, user.ID)
	// tbd: no need to return user!
	return s.Render.JSON(w, http.StatusCreated, user)
}

func login(s *Server, w http.ResponseWriter, r *http.Request) error {
	decoder := &decoders.Login{}
	if err := decoders.Decode(r, decoder); err != nil {
		return err
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {
		if isErrNoRows(err) {
			errors := decoders.Errors{
				"identifier": "No user found matching this name or email",
			}
			return errors
		}
		return err
	}

	if !user.CheckPassword(decoder.Password) {
		errors := decoders.Errors{
			"password": "Your password is invalid",
		}
		return errors
	}

	// get bookmarks & subscriptions
	if user.Bookmarks, err = s.DB.Bookmarks.SelectByUserID(user.ID); err != nil {
		return err
	}
	if user.Subscriptions, err = s.DB.Subscriptions.SelectByUserID(user.ID); err != nil {
		return err
	}

	// login user
	s.setAuthCookie(w, user.ID)

	// tbd: no need to return user!
	return s.Render.JSON(w, http.StatusOK, user)

}

func logout(s *Server, w http.ResponseWriter, r *http.Request) error {
	s.setAuthCookie(w, 0)
	return s.Render.Text(w, http.StatusOK, "Logged out")
}

func changeEmail(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	decoder := &decoders.NewEmail{}
	if err := decoders.Decode(r, decoder); err != nil {
		return err
	}

	// does this email exist?
	if exists, _ := s.DB.Users.IsEmail(decoder.Email, user.ID); exists {
		errors := decoders.Errors{
			"email": "This email address is taken",
		}
		return errors
	}

	if err := s.DB.Users.UpdateEmail(decoder.Email, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "email updated")
}

func changePassword(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	decoder := &decoders.NewPassword{}
	if err := decoders.Decode(r, decoder); err != nil {
		return err
	}

	// validate old password first

	if !user.CheckPassword(decoder.OldPassword) {
		errors := decoders.Errors{
			"oldPassword": "Your old password was not recognized",
		}
		return errors
	}
	user.SetPassword(decoder.NewPassword)

	if err := s.DB.Users.UpdatePassword(user.Password, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "password updated")
}

func deleteAccount(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	if err := s.DB.Users.DeleteUser(user.ID); err != nil {
		return err
	}
	// log user out
	s.setAuthCookie(w, 0)
	return s.Render.Text(w, http.StatusOK, "account deleted")
}

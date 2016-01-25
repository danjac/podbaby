package api

import (
	"database/sql"
	"github.com/labstack/echo"
	"math/rand"
	"net/http"
	"time"

	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	passwordChars        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomPasswordLength = 6
)

func generateRandomPassword() string {
	b := make([]byte, randomPasswordLength)
	numChars := len(passwordChars)
	for i := range b {
		b[i] = passwordChars[rand.Intn(numChars)]
	}
	return string(b)
}

func recoverPassword(c *echo.Context) error {
	// generate a temp password
	decoder := &decoders.RecoverPassword{}

	if err, ok := decoders.Decode(c, decoder); !ok {
		return err
	}

	var (
		store     = storeFromContext(c)
		userStore = store.Users()
		conn      = store.Conn()
	)
	user, err := userStore.GetByNameOrEmail(conn, decoder.Identifier)

	if err != nil {
		if err == sql.ErrNoRows {
			errors := decoders.Errors{
				"identifier": "No user found matching this email or name",
			}
			return errors.Render(c)
		}
		return err
	}

	user.SetPassword(generateRandomPassword())

	if err := userStore.UpdatePassword(conn, user.Password, user.ID); err != nil {
		return err
	}

	data := map[string]string{
		"name":         user.Name,
		"tempPassword": tempPassword,
		"host":         r.Host,
	}

	go func(c *echo.Context, to string, data map[string]string) {

		mailer := mailerFromContext(c)
		err := mailer.SendFromTemplate(
			"services@podbaby.me",
			[]string{to},
			"Your new password",
			"recover_password.tmpl",
			data,
		)
		if err != nil {
			c.Echo().Logger().Error(err)
		}

	}(c, user.Email, data)

	return c.NoContent(w, http.StatusOK)
}

func isName(c *echo.Context) error {

	var (
		err    error
		exists bool
	)

	name := c.Form("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	store := storeFromContext(c)

	if exists, err = store.Users().IsName(store.Conn(), name); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, exists)
}

func isEmail(c *echo.Context) error {
	var (
		err    error
		exists bool
		userID int64
	)

	email := c.Form("email")

	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if user, ok := userFromContextOk(c); ok {
		userID = user.ID
	}

	if exists, err = store.Users().IsEmail(store.Conn(), email, userID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, exists)

}

func signup(c *echo.Context) error {

	decoder := &decoders.Signup{}

	if err, ok := decoders.Decode(c, decoder); !ok {
		return err
	}

	var (
		store     = storeFromContext(c)
		userStore = store.Users()
		conn      = store.Conn()
	)

	errors := make(decoders.Errors)

	if exists, _ := userStore.IsEmail(conn, decoder.Email, 0); exists {
		errors["email"] = "This email address is taken"
	}

	if exists, _ := userStore.IsName(conn, decoder.Name); exists {
		errors["name"] = "This name is taken"
	}

	if len(errors) > 0 {
		return errors.Render(c)
	}

	// make new user

	user := &models.User{
		Name:  decoder.Name,
		Email: decoder.Email,
	}

	if err := user.SetPassword(decoder.Password); err != nil {
		return err
	}

	if err := userStore.Create(conn, user); err != nil {
		return err
	}

	cookieStore := cookieStoreFromContext(c)
	if err := cookieStore.Write(c.Response(), user.ID); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func login(c *echo.Context) error {

	decoder := &decoders.Login{}
	if err, ok := decoders.Decode(c, decoder); !ok {
		return err
	}

	var (
		store = storeFromContext(c)
		conn  = store.Conn()
	)

	user, err := store.Users().GetByNameOrEmail(conn, decoder.Identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			errors := decoders.Errors{
				"identifier": "No user found matching this name or email",
			}
			return errors.Render(c)
		}
		return err
	}

	if !user.CheckPassword(decoder.Password) {
		errors := decoders.Errors{
			"password": "Your password is invalid",
		}
		return errors.Render(c)
	}

	// get bookmarks & subscriptions
	if user.Bookmarks, err = store.Bookmarks().SelectByUserID(conn, user.ID); err != nil {
		return err
	}
	if user.Subscriptions, err = store.Subscriptions().SelectByUserID(conn, user.ID); err != nil {
		return err
	}
	if user.Plays, err = store.Plays().SelectByUserID(conn, user.ID); err != nil {
		return err
	}
	// login user

	cookieStore := cookieStoreFromContext(c)
	if err := cookieStore.Write(userCookieID, user.ID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)

}

func logout(c *echo.Context) error {
	cookieStore := cookieStoreFromContext(c)
	cookieStore.Write(userCookieID, 0)
	return c.NoContent(http.StatusOK)
}

func changeEmail(c *echo.Context) error {
	user := userFromContext(c)

	decoder := &decoders.NewEmail{}
	if err, ok := decoders.Decode(c, decoder); !ok {
		return err
	}

	// does this email exist?
	var (
		store     = storeFromContext(c)
		userStore = store.Users()
		conn      = store.Conn()
	)

	if exists, _ := userStore.IsEmail(conn, decoder.Email, user.ID); exists {
		errors := decoders.Errors{
			"email": "This email address is taken",
		}
		return errors.Render(c)
	}

	if err := s.DB.Users.UpdateEmail(decoder.Email, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "email updated")
}

func changePassword(c *echo.Context) error {

	user := userFromContext(c)

	decoder := &decoders.NewPassword{}
	if err, ok := decoders.Decode(c, decoder); !ok {
		return err
	}

	// validate old password first

	if !user.CheckPassword(decoder.OldPassword) {
		errors := decoders.Errors{
			"oldPassword": "Your old password was not recognized",
		}
		return errors.Render(c)
	}
	user.SetPassword(decoder.NewPassword)

	store := storeFromContext(c)

	if err := store.Users().UpdatePassword(store.Conn(), user.Password, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
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

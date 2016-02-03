package api

import (
	"database/sql"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"math/rand"
	"net/http"
	"time"

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

	var (
		v         = newValidator(c)
		store     = getStore(c)
		userStore = store.Users()
		conn      = store.Conn()
	)

	decoder := &recoverPasswordDecoder{}
	user := &models.User{}

	if ok, err := v.handleFunc(decoder, func(v *validator) error {
		if err := userStore.GetByNameOrEmail(conn, user, decoder.Identifier); err != nil {
			if err == sql.ErrNoRows {
				v.invalid("identifier", "No user found matching this email or name")
			}
			return err
		}
		return nil
	}); !ok {
		return err
	}

	tempPassword := generateRandomPassword()
	user.SetPassword(tempPassword)

	if err := userStore.UpdatePassword(conn, user.Password, user.ID); err != nil {
		return err
	}

	data := map[string]string{
		"name":         user.Name,
		"tempPassword": tempPassword,
		"host":         c.Request().Host,
	}

	go func(to string, data map[string]string) {

		if err := getMailer(c).SendFromTemplate(
			"services@podbaby.me",
			[]string{to},
			"Your new password",
			"recover_password.tmpl",
			data,
		); err != nil {
			c.Echo().Logger().Error(err)
		}

	}(user.Email, data)

	return c.NoContent(http.StatusOK)
}

func isName(c *echo.Context) error {

	var (
		err    error
		exists bool
		store  = getStore(c)
	)

	name := c.Form("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if exists, err = store.Users().IsName(store.Conn(), name); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, exists)
}

func isEmail(c *echo.Context) error {
	var (
		err    error
		exists bool
		userID int
		store  = getStore(c)
	)

	email := c.Form("email")

	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if user, _ := authenticate(c); user != nil {
		userID = user.ID
	}

	if exists, err = store.Users().IsEmail(store.Conn(), email, userID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, exists)

}

func signup(c *echo.Context) error {

	var (
		v           = newValidator(c)
		cookieStore = getCookieStore(c)
		store       = getStore(c)
		userStore   = store.Users()
		conn        = store.Conn()
	)
	decoder := &signupDecoder{}

	if ok, err := v.handleFunc(decoder, func(v *validator) error {
		if exists, _ := userStore.IsEmail(conn, decoder.Email, 0); exists {
			v.invalid("email", "This email address is taken")
		}

		if exists, _ := userStore.IsName(conn, decoder.Name); exists {
			v.invalid("name", "This name is taken")
		}
		return nil
	}); !ok {
		return err
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

	if err := cookieStore.Write(c, userCookieKey, user.ID); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func login(c *echo.Context) error {

	var (
		v           = newValidator(c)
		cookieStore = getCookieStore(c)
		store       = getStore(c)
		conn        = store.Conn()
	)

	decoder := &loginDecoder{}
	user := &models.User{}

	if ok, err := v.handleFunc(decoder, func(v *validator) error {
		if err := store.Users().GetByNameOrEmail(conn, user, decoder.Identifier); err != nil {
			if err == sql.ErrNoRows {
				v.invalid("identifier", "No user found matching this name or email")
			} else {
				return err
			}
		}

		if !user.CheckPassword(decoder.Password) {
			v.invalid("password", "Your password is invalid")
		}

		return nil

	}); !ok {
		return err
	}

	// get bookmarks & subscriptions
	if err := store.Bookmarks().SelectByUserID(conn, &user.Bookmarks, user.ID); err != nil {
		return err
	}
	if err := store.Subscriptions().SelectByUserID(conn, &user.Subscriptions, user.ID); err != nil {
		return err
	}
	if err := store.Plays().SelectByUserID(conn, &user.Plays, user.ID); err != nil {
		return err
	}
	// login user

	if err := cookieStore.Write(c, userCookieKey, user.ID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)

}

func logout(c *echo.Context) error {
	getCookieStore(c).Write(c, userCookieKey, 0)
	return c.NoContent(http.StatusOK)
}

func changeEmail(c *echo.Context) error {

	var (
		v         = newValidator(c)
		user      = getUser(c)
		store     = getStore(c)
		userStore = store.Users()
		conn      = store.Conn()
	)

	decoder := &changeEmailDecoder{}

	if ok, err := v.handleFunc(decoder, func(v *validator) error {
		if exists, _ := userStore.IsEmail(conn, decoder.Email, user.ID); exists {
			v.invalid("email", "This email address is taken")
		}
		return nil

	}); !ok {
		return err
	}

	if err := userStore.UpdateEmail(conn, decoder.Email, user.ID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func changePassword(c *echo.Context) error {

	var (
		v     = newValidator(c)
		user  = getUser(c)
		store = getStore(c)
	)

	decoder := &changePasswordDecoder{}
	if ok, err := v.handleFunc(decoder, func(v *validator) error {
		if !user.CheckPassword(decoder.OldPassword) {
			v.invalid("oldPassword", "Your old password was not recognized")
		}
		return nil

	}); !ok {
		return err
	}

	// validate old password first

	user.SetPassword(decoder.NewPassword)

	if err := store.Users().UpdatePassword(store.Conn(), user.Password, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func deleteAccount(c *echo.Context) error {

	var (
		user  = getUser(c)
		store = getStore(c)
	)

	if err := store.Users().DeleteUser(store.Conn(), user.ID); err != nil {
		return err
	}

	// log user out
	return c.NoContent(http.StatusNoContent)
}

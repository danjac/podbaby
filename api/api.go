package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/justinas/nosurf"
	"github.com/labstack/echo"
	"net/http"
	"time"

	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/store"
	"github.com/gorilla/securecookie"
)

const (
	cookieUserID          = "userid"
	userContextKey        = "user"
	storeContextKey       = "store"
	feedparserContextKey  = "feedparser"
	mailerContextKey      = "mailer"
	cookieStoreContextKey = "cookieStore"
	configContextKey      = "config"
)

type Env struct {
	*config.Config
	Store      store.Store
	Cookie     CookieStore
	Feedparser feedparser.Feedparser
	Mailer     mailer.Mailer
}

func New(env *Env) http.Handler {

	e := echo.New()

	// add all the application globals we'll need

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(storeContextKey, store)
			c.Set(mailerContextKey, mailer)
			c.Set(feedparserContextKey, feedparser)
			c.Set(configContextKey, cfg)
			c.Set(cookieStoreContextKey, cookie)
			return h(c)
		}
	})

	// static configuration
	e.Static(env.StaticURL, env.StaticDir)

	withRoutes(e)

	return nosurf.Pure(e)

}

func userFromContext(c *echo.Context) *models.User {
	return c.Get(userContextKey).(*models.User)
}

func userFromContextOk(c *echo.Context) (*models.User, bool) {
	v := c.Get(userContextKey)
	if v == nil {
		return nil, false
	}
	return c.Get(userContextKey).(*models.User), true
}

func storeFromContext(c *echo.Context) store.Store {
	return c.Get(storeContextKey).(store.Store)
}

func cookieStoreFromContext(c *echo.Context) CookieStore {
	return c.Get(cookieStoreContextKey).(CookieStore)
}

func mailerFromContext(c *echo.Context) mailer.Mailer {
	return c.Get(mailerContextKey).(mailer.Mailer)
}

func feedparserFromContext(c *echo.Context) feedparser.Feedparser {
	return c.Get(feedparserContextKey).(feedparser.Feedparser)
}

func configFromContext(c *echo.Context) *config.Config {
	return c.Get(configContextKey).(*config.Config)
}

func getInt64(c *echo.Context, name string) (int64, error) {
	value, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil {
		return value, echo.NewHTTPError(http.StatusNotFound)
	}
	return value, nil
}

func getPage(c *echo.Context) int64 {
	value := c.Form("page")
	if value == "" {
		return 1
	}
	page, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		page = 1
	}
	return page
}

func authorize() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.handlerFunc {
		return func(c *echo.Context) error {
			user, err := authenticate(c)
			if err != nil {
				return err
			}
			if user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			return h(c)
		}
	}
}

func login(c *echo.Context, userID int64) error {
	cookieStore := cookieStoreFromContext(c)
	return cookieStore.Write(userCookieID, userID)
}

func logout(c *echo.Context) error {
	cookieStore := cookieStoreFromContext(c)
	return cookieStore.Write(userCookieID, 0)
}

func authenticate(c *echo.Context) (*models.User, error) {

	// just in case this function has already been called elsewhere
	if user, ok := userFromContextOk(c); ok {
		return user, nil
	}

	cookieStore := cookieStoreFromContext(c)
	var userID int64

	if err := cookieStore.Read(cookieUserID, &userID); err != nil {
		return nil, nil
	}

	if userID == 0 {
		return nil, nil
	}

	store := storeFromContext(c)
	user, err := store.Users().GetByID(store.Conn(), userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	c.Set(userContextKey, user)
	return user, nil
}

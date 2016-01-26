package api

import (
	"database/sql"
	"github.com/justinas/nosurf"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
)

const (
	userCookieKey           = "userid"
	userContextKey          = "user"
	storeContextKey         = "store"
	feedparserContextKey    = "feedparser"
	mailerContextKey        = "mailer"
	cookieStoreContextKey   = "cookieStore"
	configContextKey        = "config"
	authenticatorContextKey = "authenticator"
)

type Env struct {
	*config.Config
	Store       store.Store
	CookieStore CookieStore
	Feedparser  feedparser.Feedparser
	Mailer      mailer.Mailer
}

type authenticator interface {
	authenticate(*echo.Context) (*models.User, error)
}

type renderer struct {
	templates *template.Template
}

// Render HTML
func (r *renderer) Render(w io.Writer, name string, data interface{}) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func New(env *Env) (http.Handler, error) {

	e := echo.New()
	e.SetDebug(true)
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Render HTML

	templates, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		return nil, err
	}
	e.SetRenderer(&renderer{templates})

	// add all the application globals we'll need

	auth := &defaultAuthenticator{}
	cookieStore := newCookieStore(env.Config)

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(storeContextKey, env.Store)
			c.Set(mailerContextKey, env.Mailer)
			c.Set(feedparserContextKey, env.Feedparser)
			c.Set(configContextKey, env.Config)
			c.Set(cookieStoreContextKey, cookieStore)
			c.Set(authenticatorContextKey, auth)
			return h(c)
		}
	})

	// catch sql no rows errors and return as a 404
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := h(c)
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			return err
		}
	})

	// static configuration
	e.Static(env.StaticURL, env.StaticDir)

	withRoutes(e)

	return nosurf.NewPure(e), nil

}

func getUser(c *echo.Context) *models.User {
	return c.Get(userContextKey).(*models.User)
}

func getUserOk(c *echo.Context) (*models.User, bool) {
	user, ok := c.Get(userContextKey).(*models.User)
	return user, ok
}

func getStore(c *echo.Context) store.Store {
	return c.Get(storeContextKey).(store.Store)
}

func getCookieStore(c *echo.Context) CookieStore {
	return c.Get(cookieStoreContextKey).(CookieStore)
}

func getMailer(c *echo.Context) mailer.Mailer {
	return c.Get(mailerContextKey).(mailer.Mailer)
}

func getFeedparser(c *echo.Context) feedparser.Feedparser {
	return c.Get(feedparserContextKey).(feedparser.Feedparser)
}

func getConfig(c *echo.Context) *config.Config {
	return c.Get(configContextKey).(*config.Config)
}

func getAuthenticator(c *echo.Context) authenticator {
	return c.Get(authenticatorContextKey).(authenticator)
}

func getIntOr404(c *echo.Context, name string) (int64, error) {
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
	return func(h echo.HandlerFunc) echo.HandlerFunc {
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

func authenticate(c *echo.Context) (*models.User, error) {
	return getAuthenticator(c).authenticate(c)
}

type defaultAuthenticator struct{}

func (a *defaultAuthenticator) authenticate(c *echo.Context) (*models.User, error) {

	// just in case this function has already been called elsewhere
	if user, ok := getUserOk(c); ok {
		return user, nil
	}

	cookieStore := getCookieStore(c)
	var userID int64

	if err := cookieStore.Read(c, userCookieKey, &userID); err != nil {
		return nil, nil
	}

	if userID == 0 {
		return nil, nil
	}

	store := getStore(c)
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

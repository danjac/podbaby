package api

import (
	"fmt"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/justinas/nosurf"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	mw "github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo/middleware"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/golang.org/x/net/http2"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/danjac/podbaby/cache"
	"github.com/danjac/podbaby/config"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/mailer"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store"
)

const (
	userSessionKey          = "userid"
	userContextKey          = "user"
	storeContextKey         = "store"
	feedparserContextKey    = "feedparser"
	cacheContextKey         = "cache"
	mailerContextKey        = "mailer"
	sessionContextKey       = "session"
	configContextKey        = "config"
	authenticatorContextKey = "authenticator"
)

type Env struct {
	*config.Config
	Cache      cache.Cache
	Store      store.Store
	Session    Session
	Feedparser feedparser.Feedparser
	Mailer     mailer.Mailer
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

func Run(env *Env) error {

	e := echo.New()
	e.SetDebug(true)
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Render HTML

	templates, err := template.ParseGlob(filepath.Join(env.TemplateDir, "*.tmpl"))
	if err != nil {
		return err
	}
	e.SetRenderer(&renderer{templates})

	// add all the application globals we'll need

	auth := &defaultAuthenticator{}
	session := newSession(env.Config)

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(storeContextKey, env.Store)
			c.Set(mailerContextKey, env.Mailer)
			c.Set(feedparserContextKey, env.Feedparser)
			c.Set(cacheContextKey, env.Cache)
			c.Set(configContextKey, env.Config)
			c.Set(sessionContextKey, session)
			c.Set(authenticatorContextKey, auth)
			return h(c)
		}
	})

	// catch sql no rows errors and return as a 404
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := h(c)
			if err == store.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			return err
		}
	})

	// static configuration
	e.Static(env.StaticURL, env.StaticDir)

	configureRoutes(e)

	// add CSRF protection
	s := e.Server(fmt.Sprintf(":%v", env.Port))
	s.Handler = nosurf.NewPure(e)
	http2.ConfigureServer(s, nil)
	return s.ListenAndServe()

}

func getUser(c *echo.Context) *models.User {
	return c.Get(userContextKey).(*models.User)
}

func getUserOk(c *echo.Context) (*models.User, bool) {
	user, ok := c.Get(userContextKey).(*models.User)
	return user, ok
}

func getCache(c *echo.Context) cache.Cache {
	return c.Get(cacheContextKey).(cache.Cache)
}

func getStore(c *echo.Context) store.Store {
	return c.Get(storeContextKey).(store.Store)
}

func getSession(c *echo.Context) Session {
	return c.Get(sessionContextKey).(Session)
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

func getIntOr404(c *echo.Context, name string) (int, error) {
	value, err := strconv.Atoi(c.Param(name))
	if err != nil {
		return value, echo.NewHTTPError(http.StatusNotFound)
	}
	return value, nil
}

func getPage(c *echo.Context) int {
	value := c.Form("page")
	if value == "" {
		return 1
	}
	page, err := strconv.Atoi(value)
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

	session := getSession(c)
	var userID int

	if err := session.Read(c, userSessionKey, &userID); err != nil {
		return nil, nil
	}

	if userID == 0 {
		return nil, nil
	}

	s := getStore(c)
	user := &models.User{}
	if err := s.Users().GetByID(s.Conn(), user, userID); err != nil {
		if err == store.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	c.Set(userContextKey, user)
	return user, nil
}

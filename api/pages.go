package api

import (
	"net/http"
	"time"

	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/justinas/nosurf"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/models"
)

func indexPage(c *echo.Context) error {

	var (
		dynamicContentURL string
		err               error
		cfg               = getConfig(c)
		store             = getStore(c)
		conn              = store.Conn()
	)

	user, err := authenticate(c)

	if err != nil {
		return err
	}

	if user != nil {
		if err = store.Bookmarks().SelectByUserID(conn, &user.Bookmarks, user.ID); err != nil {
			return err
		}
		if err = store.Subscriptions().SelectByUserID(conn, &user.Subscriptions, user.ID); err != nil {
			return err
		}
		if err = store.Plays().SelectByUserID(conn, &user.Plays, user.ID); err != nil {
			return err
		}
	}

	var categories []models.Category

	if err := getCache(c).Get("categories", time.Hour*24, &categories, func() error {
		return store.Categories().SelectAll(conn, &categories)
	}); err != nil {
		return err
	}

	csrfToken := nosurf.Token(c.Request())

	if cfg.IsDev() {
		dynamicContentURL = cfg.DynamicContentURL
	} else {
		dynamicContentURL = cfg.StaticURL
	}

	data := map[string]interface{}{
		"env":               cfg.Env,
		"dynamicContentURL": dynamicContentURL,
		"staticURL":         cfg.StaticURL,
		"googleAnalyticsID": cfg.GoogleAnalyticsID,
		"csrfToken":         csrfToken,
		"categories":        categories,
		"user":              user,
		"timestamp":         time.Now().Unix(),
	}
	return c.Render(http.StatusOK, "index.tmpl", data)
}

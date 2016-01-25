package api

import (
	"net/http"
	"time"

	"github.com/justinas/nosurf"
	"github.com/labstack/echo"
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
		if user.Bookmarks, err = store.Bookmarks().SelectByUserID(conn, user.ID); err != nil {
			return err
		}
		if user.Subscriptions, err = store.Subscriptions().SelectByUserID(conn, user.ID); err != nil {
			return err
		}
		if user.Plays, err = store.Plays().SelectByUserID(conn, user.ID); err != nil {
			return err
		}
	}

	categories, err := store.Categories().SelectAll(conn)
	if err != nil {
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

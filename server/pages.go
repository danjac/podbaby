package server

import (
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

func indexPage(s *Server, w http.ResponseWriter, r *http.Request) error {

	var (
		dynamicContentURL string
		err               error
	)

	if s.Config.IsDev() {
		dynamicContentURL = s.Config.DynamicContentURL
	} else {
		dynamicContentURL = s.Config.StaticURL
	}
	user, ok := getUser(r)
	if ok {
		if user.Bookmarks, err = s.DB.Bookmarks.SelectByUserID(user.ID); err != nil {
			return err
		}
		if user.Subscriptions, err = s.DB.Subscriptions.SelectByUserID(user.ID); err != nil {
			return err
		}
		if user.Plays, err = s.DB.Plays.SelectByUserID(user.ID); err != nil {
			return err
		}
	}

	categories, err := s.DB.Categories.SelectAll()
	if err != nil {
		return err
	}

	csrfToken := nosurf.Token(r)
	ctx := map[string]interface{}{
		"env":               s.Config.Env,
		"dynamicContentURL": dynamicContentURL,
		"staticURL":         s.Config.StaticURL,
		"googleAnalyticsID": s.Config.GoogleAnalyticsID,
		"csrfToken":         csrfToken,
		"categories":        categories,
		"user":              user,
		"timestamp":         time.Now().Unix(),
	}
	return s.Render.HTML(w, http.StatusOK, "index", ctx)
}

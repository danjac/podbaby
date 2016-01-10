package server

import (
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) {

	var dynamicContentURL string
	if s.Config.Env == "dev" {
		dynamicContentURL = s.Config.DynamicContentURL
	} else {
		dynamicContentURL = s.Config.StaticURL
	}
	user, err := s.getUserFromCookie(r)
	if err == nil {
		if user.Bookmarks, err = s.DB.Bookmarks.SelectByUserID(user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
		if user.Subscriptions, err = s.DB.Subscriptions.SelectByUserID(user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
	}
	csrfToken := nosurf.Token(r)
	ctx := map[string]interface{}{
		"dynamicContentURL": dynamicContentURL,
		"staticURL":         s.Config.StaticURL,
		"googleAnalyticsID": s.Config.GoogleAnalyticsID,
		"csrfToken":         csrfToken,
		"user":              user,
		"timestamp":         time.Now().Unix(),
	}
	s.Render.HTML(w, http.StatusOK, "index", ctx)
}

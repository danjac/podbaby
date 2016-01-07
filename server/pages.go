package server

import (
	"net/http"

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
		"csrfToken":         csrfToken,
		"user":              user,
	}
	s.Render.HTML(w, http.StatusOK, "index", ctx)
}

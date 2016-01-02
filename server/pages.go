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
	user, _ := s.getUserFromCookie(r)
	csrfToken := nosurf.Token(r)
	ctx := map[string]interface{}{
		"dynamicContentURL": dynamicContentURL,
		"staticURL":         s.Config.StaticURL,
		"csrfToken":         csrfToken,
		"user":              user,
	}
	s.Render.HTML(w, http.StatusOK, "index", ctx)
}

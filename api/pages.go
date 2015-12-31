package api

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) {
	user, _ := s.getUserFromCookie(r)
	csrfToken := nosurf.Token(r)
	ctx := map[string]interface{}{
		"dynamicContentURL": s.Config.DynamicContentURL,
		"staticURL":         s.Config.StaticURL,
		"csrfToken":         csrfToken,
		"user":              user,
	}
	s.Render.HTML(w, http.StatusOK, "index", ctx)
}

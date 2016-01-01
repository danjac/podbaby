package server

import (
	"net/http"
	"strings"

	"github.com/danjac/podbaby/models"
)

func (s *Server) search(w http.ResponseWriter, r *http.Request) {

	user, _ := getUser(r)
	query := strings.Trim(r.FormValue("q"), " ")

	result := &models.SearchResult{}

	if query != "" {
		var err error
		if result.Channels, err = s.DB.Channels.Search(query, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
		if result.Podcasts, err = s.DB.Podcasts.Search(query, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
	}

	s.Render.JSON(w, http.StatusOK, result)
}

package server

import (
	"net/http"
)

func (s *Server) getLatestPodcasts(w http.ResponseWriter, r *http.Request) {
	podcasts, err := s.DB.Podcasts.SelectAll()
	if err != nil {
		s.Abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, podcasts)

}

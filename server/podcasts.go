package server

import (
	"net/http"
)

func (s *Server) getLatestPodcasts(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	s.Log.Info("current user:" + user.Name)
	podcasts, err := s.DB.Podcasts.SelectAll()
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, podcasts)

}

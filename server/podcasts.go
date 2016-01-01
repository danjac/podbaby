package server

import "net/http"

func (s *Server) getLatestPodcasts(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	result, err := s.DB.Podcasts.SelectSubscribed(user.ID, getPage(r))
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, result)
}

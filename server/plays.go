package server

import "net/http"

func (s *Server) addPlay(w http.ResponseWriter, r *http.Request) {
	podcastID, _ := getInt64(r, "id")
	user, _ := getUser(r)
	if err := s.DB.Plays.Create(podcastID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusCreated, "played")
}

func (s *Server) getPlays(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	result, err := s.DB.Podcasts.SelectPlayed(user.ID, getPage(r))
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, result)
}

func (s *Server) deleteAllPlays(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	if err := s.DB.Plays.DeleteAll(user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "all plays removed")
}

package server

import "net/http"

func addPlay(s *Server, w http.ResponseWriter, r *http.Request) error {
	podcastID, _ := getID(r)
	user, _ := getUser(r)
	if err := s.DB.Plays.Create(podcastID, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusCreated, "played")
}

func getPlays(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	result, err := s.DB.Podcasts.SelectPlayed(user.ID, getPage(r))
	if err != nil {
		return err
	}
	return s.Render.JSON(w, http.StatusOK, result)
}

func deleteAllPlays(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	if err := s.DB.Plays.DeleteAll(user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "all plays removed")
}

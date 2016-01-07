package server

import (
	"github.com/danjac/podbaby/models"
	"net/http"
)

func (s *Server) getPodcast(w http.ResponseWriter, r *http.Request) {

	podcastID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}
	podcast, err := s.DB.Podcasts.GetByID(podcastID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, podcast)
}

func (s *Server) getLatestPodcasts(w http.ResponseWriter, r *http.Request) {
	user, err := s.getUserFromCookie(r)

	var (
		result *models.PodcastList
	)

	if err == nil { // user authenticated
		result, err = s.DB.Podcasts.SelectSubscribed(user.ID, getPage(r))
	} else {
		result, err = s.DB.Podcasts.SelectAll(getPage(r))
	}

	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, result)
}

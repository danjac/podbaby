package server

import (
	"github.com/danjac/podbaby/models"
	"net/http"
)

func getPodcast(s *Server, w http.ResponseWriter, r *http.Request) error {

	podcastID, _ := getID(r)
	podcast, err := s.DB.Podcasts.GetByID(podcastID)
	if err != nil {
		return err
	}
	return s.Render.JSON(w, http.StatusOK, podcast)
}

func getLatestPodcasts(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, ok := getUser(r)

	var (
		result *models.PodcastList
		err    error
	)

	if ok { // user authenticated
		result, err = s.DB.Podcasts.SelectSubscribed(user.ID, getPage(r))
	} else {
		result, err = s.DB.Podcasts.SelectAll(getPage(r))
	}

	if err != nil {
		return err
	}
	return s.Render.JSON(w, http.StatusOK, result)
}

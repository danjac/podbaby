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

	var (
		result *models.PodcastList
		err    error
	)

	user, ok := getUser(r)
	page := getPage(r)

	if ok { // user authenticated
		result, err = s.DB.Podcasts.SelectSubscribed(user.ID, page)
	} else {
		result, err = s.DB.Podcasts.SelectAll(page)
	}

	if err != nil {
		return err
	}
	return s.Render.JSON(w, http.StatusOK, result)
}

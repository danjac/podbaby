package server

import (
	"github.com/danjac/podbaby/models"
	"net/http"
)

func (s *Server) getLatestPodcasts(w http.ResponseWriter, r *http.Request) {
	user, ok := getUser(r)

	var (
		err    error
		result *models.PodcastList
	)

	if ok {
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

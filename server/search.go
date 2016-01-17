package server

import (
	"net/http"
	"strings"

	"github.com/danjac/podbaby/models"
)

func searchAll(s *Server, w http.ResponseWriter, r *http.Request) error {

	query := strings.Trim(r.FormValue("q"), " ")

	result := &models.SearchResult{}

	if query != "" {
		var err error
		if result.Channels, err = s.DB.Channels.Search(query); err != nil {
			return err
		}
		if result.Podcasts, err = s.DB.Podcasts.Search(query); err != nil {
			return err
		}
	}

	return s.Render.JSON(w, http.StatusOK, result)
}

func searchBookmarks(s *Server, w http.ResponseWriter, r *http.Request) error {

	user, _ := getUser(r)
	query := strings.Trim(r.FormValue("q"), " ")

	var podcasts []models.Podcast
	var err error

	if query != "" {
		if podcasts, err = s.DB.Podcasts.SearchBookmarked(query, user.ID); err != nil {
			return err
		}
	}

	return s.Render.JSON(w, http.StatusOK, podcasts)
}

func searchChannel(s *Server, w http.ResponseWriter, r *http.Request) error {

	query := strings.Trim(r.FormValue("q"), " ")

	channelID, err := getID(r)
	if err != nil {
		return err
	}

	var podcasts []models.Podcast

	if query != "" {
		if podcasts, err = s.DB.Podcasts.SearchByChannelID(query, channelID); err != nil {
			return err
		}
	}

	return s.Render.JSON(w, http.StatusOK, podcasts)
}

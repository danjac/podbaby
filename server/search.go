package server

import (
	"net/http"
	"strings"

	"github.com/danjac/podbaby/models"
)

func (s *Server) search(w http.ResponseWriter, r *http.Request) {

	user, _ := getUser(r)
	query := strings.Trim(r.FormValue("q"), " ")

	result := &models.SearchResult{}

	if query != "" {
		var err error
		if result.Channels, err = s.DB.Channels.Search(query, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
		if result.Podcasts, err = s.DB.Podcasts.Search(query, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
	}

	s.Render.JSON(w, http.StatusOK, result)
}

func (s *Server) searchBookmarks(w http.ResponseWriter, r *http.Request) {

	user, _ := getUser(r)
	query := strings.Trim(r.FormValue("q"), " ")

	var podcasts []models.Podcast
	var err error

	if query != "" {
		if podcasts, err = s.DB.Podcasts.SearchBookmarked(query, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
	}

	s.Render.JSON(w, http.StatusOK, podcasts)
}

func (s *Server) searchChannel(w http.ResponseWriter, r *http.Request) {

	user, _ := getUser(r)
	query := strings.Trim(r.FormValue("q"), " ")

	channelID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}

	var podcasts []models.Podcast

	if query != "" {
		if podcasts, err = s.DB.Podcasts.SearchByChannelID(query, channelID, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
	}

	s.Render.JSON(w, http.StatusOK, podcasts)
}

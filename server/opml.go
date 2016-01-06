package server

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"net/http"
)

func (s *Server) getOPML(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channels, err := s.DB.Channels.SelectSubscribed(user.ID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	opml := &models.OPML{
		Version: "1.0",
		Title:   fmt.Sprintf("Podcast subscriptions for %s", user.Name),
	}
	for _, channel := range channels {
		var website string
		if channel.Website.Valid {
			website = channel.Website.String
		}
		opml.Outlines = append(opml.Outlines, &models.Outline{
			Type:    "rss",
			Title:   channel.Title,
			Text:    channel.Title,
			URL:     channel.URL,
			HtmlURL: website,
		},
		)
	}
	s.Render.XML(w, http.StatusOK, opml)
}

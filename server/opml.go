package server

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"net/http"
)

func getOPML(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	channels, err := s.DB.Channels.SelectSubscribed(user.ID)
	if err != nil {
		return err
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
	return s.Render.XML(w, http.StatusOK, opml)
}

package api

import (
	"fmt"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/models"
	"net/http"
)

func getOPML(c *echo.Context) error {

	var (
		user     = getUser(c)
		store    = getStore(c)
		channels []models.Channel
	)

	if err := store.Channels().SelectSubscribed(store.Conn(), &channels, user.ID); err != nil {
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
			HTMLURL: website,
		},
		)
	}
	return c.XML(http.StatusOK, opml)
}

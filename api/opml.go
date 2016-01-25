package api

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
	"net/http"
)

func getOPML(c *echo.Context) error {

	var (
		user  = getUser(c)
		store = getStore(c)
	)

	channels, err := store.Channels().SelectSubscribed(store.Conn(), user.ID)
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
	return c.XML(http.StatusOK, opml)
}

package api

import (
	"github.com/labstack/echo"
	"net/http"
	"strings"

	"github.com/danjac/podbaby/models"
)

func searchAll(c *echo.Context) error {

	var (
		store = getStore(c)
		conn  = store.Conn()
	)

	query := strings.Trim(c.Form("q"), " ")

	result := &models.SearchResult{}

	if query != "" {
		var err error
		if result.Channels, err = store.Channels().Search(conn, query); err != nil {
			return err
		}
		if result.Podcasts, err = store.Podcasts().Search(conn, query); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, result)
}

func searchBookmarks(c *echo.Context) error {

	var (
		store = getStore(c)
		user  = getUser(c)
	)
	query := strings.Trim(c.Form("q"), " ")

	var podcasts []models.Podcast
	var err error

	if query != "" {
		if podcasts, err = store.Podcasts().SearchBookmarked(store.Conn(), query, user.ID); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, podcasts)
}

func searchChannel(c *echo.Context) error {

	channelID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	query := strings.Trim(c.Form("q"), " ")
	store := getStore(c)

	var podcasts []models.Podcast

	if query != "" {
		if podcasts, err = store.Podcasts().SearchByChannelID(store.Conn(), query, channelID); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, podcasts)
}

package api

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strings"
	"time"

	"github.com/danjac/podbaby/models"
)

func searchAll(c *echo.Context) error {

	var (
		cache   = getCache(c)
		query   = strings.ToLower(strings.Trim(c.Form("q"), " "))
		result  = &models.SearchResult{}
		key     = fmt.Sprintf("search:all:%v", query)
		timeout = time.Minute * 30
	)

	if query != "" {
		err := cache.Get(key, timeout, result, func() error {

			var (
				store = getStore(c)
				conn  = store.Conn()
			)

			var err error
			if result.Channels, err = store.Channels().Search(conn, query); err != nil {
				return err
			}
			if result.Podcasts, err = store.Podcasts().Search(conn, query); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
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

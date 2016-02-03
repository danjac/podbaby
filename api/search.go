package api

import (
	"fmt"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
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
		if err := cache.Get(key, timeout, result, func() error {

			var (
				s    = getStore(c)
				conn = s.Conn()
			)

			if err := s.Channels().Search(conn, &result.Channels, query); err != nil {
				return err
			}
			if err := s.Podcasts().Search(conn, &result.Podcasts, query); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, result)
}

func searchBookmarks(c *echo.Context) error {

	var (
		s    = getStore(c)
		user = getUser(c)
	)
	query := strings.ToLower(strings.Trim(c.Form("q"), " "))

	var podcasts []models.Podcast

	if query != "" {
		if err := s.Podcasts().SearchBookmarked(s.Conn(), &podcasts, query, user.ID); err != nil {
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

	query := strings.ToLower(strings.Trim(c.Form("q"), " "))
	cacheKey := fmt.Sprintf("search:channel:%v:%v", channelID, query)

	var podcasts []models.Podcast

	if query != "" {
		if err := getCache(c).Get(cacheKey, time.Minute*30, &podcasts, func() error {
			s := getStore(c)
			return s.Podcasts().SearchByChannelID(s.Conn(), &podcasts, query, channelID)
		}); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, podcasts)
}

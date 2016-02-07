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
		cache      = getCache(c)
		page       = getPage(c)
		query      = strings.ToLower(strings.TrimSpace(c.Form("q")))
		searchType = c.Form("t")
		result     = models.NewSearchResult(page)
		timeout    = time.Minute * 30
	)

	if searchType != "podcasts" && searchType != "channels" {
		searchType = "podcasts"
	}

	key := fmt.Sprintf("search:all:type:%v:page:%v:%v", searchType, page, query)

	if query != "" {
		if err := cache.Get(key, timeout, result, func() error {

			var (
				s    = getStore(c)
				conn = s.Conn()
			)

			if searchType == "channels" {
				if err := s.Channels().Search(conn, &result.Channels, query); err != nil {
					return err
				}
			} else {
				if err := s.Podcasts().Search(conn, result.Podcasts, query, page); err != nil {
					return err
				}
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
	query := strings.ToLower(strings.TrimSpace(c.Form("q")))

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

	query := strings.ToLower(strings.TrimSpace(c.Form("q")))
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

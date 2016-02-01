package api

import (
	"fmt"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/models"
	"net/http"
	"time"
)

func getPodcast(c *echo.Context) error {

	podcastID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}
	var (
		cache   = getCache(c)
		key     = fmt.Sprintf("podcast:%v", podcastID)
		timeout = time.Hour * 24
	)
	podcast := &models.Podcast{}

	if err := cache.Get(key, timeout, podcast, func() error {
		store := getStore(c)
		return store.Podcasts().GetByID(store.Conn(), podcast, podcastID)
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, podcast)
}

func getLatestPodcasts(c *echo.Context) error {

	var (
		err          error
		result       = &models.PodcastList{}
		page         = getPage(c)
		store        = getStore(c)
		conn         = store.Conn()
		podcastStore = store.Podcasts()
	)

	user, err := authenticate(c)

	if user != nil { // user authenticated
		err = podcastStore.SelectSubscribed(conn, result, user.ID, page)
	} else {
		err = getCache(c).Get(fmt.Sprintf("latest-podcasts:%v", page), time.Minute*30, result, func() error {
			return podcastStore.SelectAll(conn, result, page)
		})
	}

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

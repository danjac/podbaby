package api

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func getPodcast(c *echo.Context) error {

	podcastID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	store := getStore(c)

	podcast, err := store.Podcasts().GetByID(store.Conn(), podcastID)
	if err != nil {
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
		result, err = podcastStore.SelectSubscribed(conn, user.ID, page)
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

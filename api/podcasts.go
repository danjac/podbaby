package api

import (
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
	"net/http"
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
		result       *models.PodcastList
		page         = getPage(c)
		store        = getStore(c)
		conn         = store.Conn()
		podcastStore = store.Podcasts()
	)

	user, ok := getUserOk(c)

	if ok { // user authenticated
		result, err = podcastStore.SelectSubscribed(conn, user.ID, page)
	} else {
		result, err = podcastStore.SelectAll(conn, page)
	}

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

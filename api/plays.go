package api

import (
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
	"net/http"
)

func addPlay(c *echo.Context) error {

	var (
		user  = getUser(c)
		store = getStore(c)
	)

	podcastID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	if err := store.Plays().Create(store.Conn(), podcastID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func getPlays(c *echo.Context) error {

	var (
		user   = getUser(c)
		store  = getStore(c)
		result = &models.PodcastList{}
	)

	if err := store.Podcasts().SelectPlayed(store.Conn(), result, user.ID, getPage(c)); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func deleteAllPlays(c *echo.Context) error {
	var (
		user  = getUser(c)
		store = getStore(c)
	)

	if err := store.Plays().DeleteAll(store.Conn(), user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

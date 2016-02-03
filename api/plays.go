package api

import (
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/models"
	"net/http"
)

func addPlay(c *echo.Context) error {

	var (
		user = getUser(c)
		s    = getStore(c)
	)

	podcastID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	if err := s.Plays().Create(s.Conn(), podcastID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func getPlays(c *echo.Context) error {

	var (
		user   = getUser(c)
		s      = getStore(c)
		result = &models.PodcastList{}
	)

	if err := s.Podcasts().SelectPlayed(s.Conn(), result, user.ID, getPage(c)); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func deleteAllPlays(c *echo.Context) error {
	var (
		user = getUser(c)
		s    = getStore(c)
	)

	if err := s.Plays().DeleteAll(s.Conn(), user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

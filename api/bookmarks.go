package api

import (
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/models"
	"net/http"
)

func getBookmarks(c *echo.Context) error {

	var (
		user   = getUser(c)
		s      = getStore(c)
		result = &models.PodcastList{}
	)

	if err := s.Podcasts().SelectBookmarked(
		s.Conn(),
		result,
		user.ID,
		getPage(c)); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func addBookmark(c *echo.Context) error {
	podcastID, err := getIntOr404(c, "id")

	if err != nil {
		return err
	}

	var (
		user = getUser(c)
		s    = getStore(c)
	)

	if err := s.Bookmarks().Create(s.Conn(), podcastID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func removeBookmark(c *echo.Context) error {
	podcastID, err := getIntOr404(c, "id")

	if err != nil {
		return err
	}

	var (
		user = getUser(c)
		s    = getStore(c)
	)

	if err := s.Bookmarks().Delete(s.Conn(), podcastID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

package api

import (
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
	"net/http"
)

func getBookmarks(c *echo.Context) error {

	var (
		user   = getUser(c)
		store  = getStore(c)
		result = &models.PodcastList{}
	)

	err := store.Podcasts().SelectBookmarked(
		store.Conn(),
		result,
		user.ID,
		getPage(c))

	if err != nil {
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
		user  = getUser(c)
		store = getStore(c)
	)

	if err := store.Bookmarks().Create(store.Conn(), podcastID, user.ID); err != nil {
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
		user  = getUser(c)
		store = getStore(c)
	)

	if err := store.Bookmarks().Delete(store.Conn(), podcastID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

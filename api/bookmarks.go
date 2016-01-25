package api

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func getBookmarks(c *echo.Context) error {

	var (
		user  = getUser(c)
		store = getStore(c)
	)

	fmt.Println(c.Path())
	fmt.Println(getIntOr404(c, "id"))

	result, err := store.Podcasts().SelectBookmarked(
		store.Conn(),
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

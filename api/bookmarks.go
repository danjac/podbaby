package api

import "net/http"

func getBookmarks(c *echo.Context) error {

	user := userFromContext(c)
	store := storeFromContext(c)

	result, err := store.Podcasts().SelectBookmarked(
		store.Conn,
		user.ID,
		getPage(c))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func addBookmark(c *echo.Context) error {
	user := userFromContext(c)
	podcastID, err := getInt64(c, "id")

	if err != nil {
		return err
	}

	store := storeFromContext(c)

	if err := store.Bookmarks().Create(store.Conn, podcastID, user.ID); err != nil {
		return err
	}
	c.NoContent(http.StatusCreated)
}

func removeBookmark(c *echo.Context) error {
	user := userFromContext(c)
	podcastID, err := getInt64(c, "id")

	if err != nil {
		return err
	}

	store := storeFromContext(c)

	if err := store.Bookmarks().Delete(store.Conn, podcastID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

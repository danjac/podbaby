package api

import (
	"github.com/labstack/echo"
	"net/http"
)

func subscribe(c *echo.Context) error {

	var (
		user  = getUser(c)
		store = getStore(c)
	)

	channelID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	if err := store.Subscriptions().Create(store.Conn(), channelID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func unsubscribe(c *echo.Context) error {

	var (
		user  = getUser(c)
		store = getStore(c)
	)

	channelID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	if err := store.Subscriptions().Delete(store.Conn(), channelID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

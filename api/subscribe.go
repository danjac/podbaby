package api

import (
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"net/http"
)

func subscribe(c *echo.Context) error {

	var (
		user = getUser(c)
		s    = getStore(c)
	)

	channelID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	if err := s.Subscriptions().Create(s.Conn(), channelID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func unsubscribe(c *echo.Context) error {

	var (
		user = getUser(c)
		s    = getStore(c)
	)

	channelID, err := getIntOr404(c, "id")
	if err != nil {
		return err
	}

	if err := s.Subscriptions().Delete(s.Conn(), channelID, user.ID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

package api

import (
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
)

type fakeAuthenticator struct {
	user *models.User
}

func (a *fakeAuthenticator) authenticate(c *echo.Context) (*models.User, error) {
	return a.user, nil
}

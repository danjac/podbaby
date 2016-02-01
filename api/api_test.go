package api

import (
	"github.com/danjac/podbaby/cache"
	"github.com/danjac/podbaby/models"
	"github.com/labstack/echo"
	"time"
)

type fakeAuthenticator struct {
	user *models.User
}

func (a *fakeAuthenticator) authenticate(c *echo.Context) (*models.User, error) {
	return a.user, nil
}

type fakeCache struct{}

func (c *fakeCache) Delete(string) error { return nil }

func (c *fakeCache) Get(_ string, _ time.Duration, _ interface{}, fn cache.Setter) error {
	return fn()
}

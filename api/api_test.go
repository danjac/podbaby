package api

import (
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
	"github.com/danjac/podbaby/cache"
	"github.com/danjac/podbaby/models"
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

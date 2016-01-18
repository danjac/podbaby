package server

import (
	"errors"
	"github.com/danjac/podbaby/models"
	"net/http"
)

var errMockDBError = errors.New("Fake DB error")

func mockGetVars(vars map[string]string) varsGetter {
	return func(r *http.Request) map[string]string {
		return vars
	}
}

func mockGetContext(ctx map[string]interface{}) contextGetter {
	return func(r *http.Request, key string) (interface{}, bool) {
		result, ok := ctx[key]
		return result, ok
	}
}

// set up context map with a specific user
func mockGetContextWithUser(user *models.User) contextGetter {
	ctx := map[string]interface{}{
		userKey: user,
	}
	return mockGetContext(ctx)
}

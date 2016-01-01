package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/danjac/podbaby/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type contextGetter func(*http.Request, string) (interface{}, bool)
type varsGetter func(*http.Request) map[string]string

var getVars varsGetter
var getContext contextGetter

// authentication methods

func init() {
	// wraps mux.Vars for easier testing

	getVars = func(r *http.Request) map[string]string {
		return mux.Vars(r)
	}

	// wraps context

	getContext = func(r *http.Request, key string) (interface{}, bool) {
		return context.GetOk(r, key)
	}
}

func getInt64(r *http.Request, name string) (int64, error) {
	badRequest := HTTPError{http.StatusBadRequest, errors.New("Invalid parameter for " + name)}
	value, ok := getVars(r)[name]
	if !ok {
		return 0, badRequest
	}
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, badRequest
	}
	return intval, nil
}

func getPage(r *http.Request) int64 {
	value := r.FormValue("page")
	if value == "" {
		return 1
	}
	page, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		page = 1
	}
	return page
}

func getUser(r *http.Request) (*models.User, bool) {
	val, ok := getContext(r, userKey)
	if !ok {
		return nil, false
	}
	return val.(*models.User), true
}

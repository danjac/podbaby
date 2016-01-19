package server

import (
	"net/http"
	"strconv"

	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/models"
)

func isErrNoRows(err error) bool {
	if dbErr, ok := err.(database.DBError); ok {
		return dbErr.IsNoRows()
	}
	return false
}

func getID(r *http.Request) (int64, error) {

	value, ok := getVars(r)["id"]
	if !ok {
		return 0, errBadRequest
	}
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errBadRequest
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

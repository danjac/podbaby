package server

import (
	"database/sql"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Error interface {
	error
	Status() int
}

type HTTPError struct {
	Code int
	Err  error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

func (e HTTPError) Status() int {
	return e.Code
}

func (s *Server) abort(w http.ResponseWriter, r *http.Request, err error) {
	logger := s.Log.WithFields(logrus.Fields{
		"URL":    r.URL,
		"Method": r.Method,
		"Error":  err,
	})
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger.Debug("Not found:" + err.Error())
		return
	}

	var msg string

	switch e := err.(error).(type) {
	case Error:
		msg = "HTTP Error"
		http.Error(w, e.Error(), e.Status())
	default:
		msg = "Internal Server Error"
		http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
	}
	logger.Error(msg)
}

package store

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoRows = errors.New("No rows found")

type DBError interface {
	error
	Query() string
}

type dbError struct {
	err   error
	query string
}

func (e *dbError) Query() string {
	return e.query
}

func (e *dbError) Error() string {
	return fmt.Sprintf("Err: %v\r\nSQL: %v", e.err.Error(), e.Query())
}

func handleError(err error, query string) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return ErrNoRows
	}
	return &dbError{err, query}
}

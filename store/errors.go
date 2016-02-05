package store

import (
	"database/sql"
	"errors"
	"fmt"
)

// ErrNoRows is returned if expected result is empty
var ErrNoRows = errors.New("No rows found")

// DBError returns basic error plus original query
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

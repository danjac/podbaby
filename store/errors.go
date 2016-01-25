package store

import (
	"fmt"
)

type DBError interface {
	error
	Query() string
}

type sqlError struct {
	err error
	sql string
}

func (e sqlError) Error() string {
	return fmt.Sprintf("%s:%s", e.err.Error(), e.Query())
}

func (e sqlError) Query() string {
	return e.sql
}

func dbErr(err error, query string) error {
	if err == nil {
		return nil
	}
	return sqlError{err, query}
}

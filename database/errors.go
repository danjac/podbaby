package database

import "database/sql"

type DBError interface {
	error
	Query() string
	IsNoRows() bool
}

type sqlError struct {
	err error
	sql string
}

func (e sqlError) Error() string {
	return e.err.Error()
}

func (e sqlError) Query() string {
	return e.sql
}

func (e sqlError) IsNoRows() bool {
	return e.err == sql.ErrNoRows
}

func sqlErr(err error, sql string) error {
	if err == nil {
		return nil
	}
	return sqlError{err, sql}
}

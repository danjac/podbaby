package database

import (
	"github.com/jmoiron/sqlx"
)

type PlayWriter interface {
	Create(int64, int64) error
	DeleteAll(int64) error
}

type PlayDB struct {
	PlayWriter
}

func newPlayDB(db sqlx.Ext) *PlayDB {
	return &PlayDB{
		PlayWriter: &PlayDBWriter{db},
	}
}

type PlayDBWriter struct {
	sqlx.Ext
}

func (db *PlayDBWriter) Create(podcastID, userID int64) error {
	q := "SELECT add_play($1, $2)"
	_, err := db.Exec(q, podcastID, userID)
	return sqlErr(err, q)
}

func (db *PlayDBWriter) DeleteAll(userID int64) error {
	q := "DELETE FROM plays WHERE user_id=$1"
	_, err := db.Exec(q, userID)
	return sqlErr(err, q)
}

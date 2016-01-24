package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type PlayReader interface {
	SelectByUserID(DataHandler, int64) ([]models.Play, error)
}

type PlayWriter interface {
	Create(DataHandler, int64, int64) error
	DeleteAll(DataHandler, int64) error
}

type PlayDB interface {
	PlayReader
	PlayWriter
}

type PlaySqlDB struct {
	PlayReader
	PlayWriter
}

func newPlayDB() PlayDB {
	return &PlaySqlDB{
		PlayReader: &PlaySqlReader{},
		PlayWriter: &PlaySqlWriter{},
	}
}

type PlaySqlReader struct{}

func (r *PlaySqlReader) SelectByUserID(dh DataHandler, userID int64) ([]models.Play, error) {
	q := "SELECT podcast_id, created_at FROM plays WHERE user_id=$1"
	var plays []models.Play
	err := sqlx.Select(dh, &plays, q, userID)
	return plays, dbErr(err, q)
}

type PlaySqlWriter struct{}

func (w *PlaySqlWriter) Create(dh DataHandler, podcastID, userID int64) error {
	q := "SELECT add_play($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return dbErr(err, q)
}

func (w *PlaySqlWriter) DeleteAll(dh DataHandler, userID int64) error {
	q := "DELETE FROM plays WHERE user_id=$1"
	_, err := dh.Exec(q, userID)
	return dbErr(err, q)
}

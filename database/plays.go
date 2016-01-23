package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type PlayReader interface {
	SelectByUserID(int64) ([]models.Play, error)
}

type PlayWriter interface {
	Create(int64, int64) error
	DeleteAll(int64) error
}

type PlayDB struct {
	PlayReader
	PlayWriter
}

func newPlayDB(db *sqlx.DB) *PlayDB {
	return &PlayDB{
		PlayReader: &PlayDBReader{db},
		PlayWriter: &PlayDBWriter{db},
	}
}

type PlayDBReader struct {
	*sqlx.DB
}

func (db *PlayDBReader) SelectByUserID(userID int64) ([]models.Play, error) {
	q := "SELECT podcast_id, created_at FROM plays WHERE user_id=$1"
	var plays []models.Play
	err := sqlx.Select(db, &plays, q, userID)
	return plays, dbErr(err, q)
}

type PlayDBWriter struct {
	*sqlx.DB
}

func (db *PlayDBWriter) Create(podcastID, userID int64) error {
	q := "SELECT add_play($1, $2)"
	_, err := db.Exec(q, podcastID, userID)
	return dbErr(err, q)
}

func (db *PlayDBWriter) DeleteAll(userID int64) error {
	q := "DELETE FROM plays WHERE user_id=$1"
	_, err := db.Exec(q, userID)
	return dbErr(err, q)
}

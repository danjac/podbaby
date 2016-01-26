package store

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

type PlayStore interface {
	PlayReader
	PlayWriter
}

type playSqlStore struct {
	PlayReader
	PlayWriter
}

func newPlayStore() PlayStore {
	return &playSqlStore{
		PlayReader: &playSqlReader{},
		PlayWriter: &playSqlWriter{},
	}
}

type playSqlReader struct{}

func (r *playSqlReader) SelectByUserID(dh DataHandler, userID int64) ([]models.Play, error) {
	q := "SELECT podcast_id, created_at FROM plays WHERE user_id=$1"
	var plays []models.Play
	err := sqlx.Select(dh, &plays, q, userID)
	return plays, err
}

type playSqlWriter struct{}

func (w *playSqlWriter) Create(dh DataHandler, podcastID, userID int64) error {
	q := "SELECT add_play($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return err
}

func (w *playSqlWriter) DeleteAll(dh DataHandler, userID int64) error {
	q := "DELETE FROM plays WHERE user_id=$1"
	_, err := dh.Exec(q, userID)
	return err
}

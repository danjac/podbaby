package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type PlayReader interface {
	SelectByUserID(DataHandler, *[]models.Play, int) error
}

type PlayWriter interface {
	Create(DataHandler, int, int) error
	DeleteAll(DataHandler, int) error
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

func (r *playSqlReader) SelectByUserID(dh DataHandler, plays *[]models.Play, userID int) error {
	q := "SELECT podcast_id, created_at FROM plays WHERE user_id=$1"
	return sqlx.Select(dh, plays, q, userID)
}

type playSqlWriter struct{}

func (w *playSqlWriter) Create(dh DataHandler, podcastID, userID int) error {
	q := "SELECT add_play($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return err
}

func (w *playSqlWriter) DeleteAll(dh DataHandler, userID int) error {
	q := "DELETE FROM plays WHERE user_id=$1"
	_, err := dh.Exec(q, userID)
	return err
}

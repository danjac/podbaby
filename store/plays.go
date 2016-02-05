package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

// PlayReader handles reads from play data store
type PlayReader interface {
	SelectByUserID(DataHandler, *[]models.Play, int) error
}

// PlayWriter handles writes to play data store
type PlayWriter interface {
	Create(DataHandler, int, int) error
	DeleteAll(DataHandler, int) error
}

// PlayStore handles interactions with play data store
type PlayStore interface {
	PlayReader
	PlayWriter
}

type playSQLStore struct {
	PlayReader
	PlayWriter
}

func newPlayStore() PlayStore {
	return &playSQLStore{
		PlayReader: &playSQLReader{},
		PlayWriter: &playSQLWriter{},
	}
}

type playSQLReader struct{}

func (r *playSQLReader) SelectByUserID(dh DataHandler, plays *[]models.Play, userID int) error {
	q := "SELECT podcast_id, created_at FROM plays WHERE user_id=$1"
	return handleError(sqlx.Select(dh, plays, q, userID), q)
}

type playSQLWriter struct{}

func (w *playSQLWriter) Create(dh DataHandler, podcastID, userID int) error {
	q := "SELECT add_play($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return handleError(err, q)
}

func (w *playSQLWriter) DeleteAll(dh DataHandler, userID int) error {
	q := "DELETE FROM plays WHERE user_id=$1"
	_, err := dh.Exec(q, userID)
	return handleError(err, q)
}

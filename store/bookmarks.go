package store

import (
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type BookmarkWriter interface {
	Create(DataHandler, int, int) error
	Delete(DataHandler, int, int) error
}

type BookmarkReader interface {
	SelectByUserID(DataHandler, *[]int, int) error
}

type BookmarkStore interface {
	BookmarkReader
	BookmarkWriter
}

type bookmarkSqlStore struct {
	BookmarkReader
	BookmarkWriter
}

func newBookmarkStore() BookmarkStore {
	return &bookmarkSqlStore{
		BookmarkReader: &bookmarkSqlReader{},
		BookmarkWriter: &bookmarkSqlWriter{},
	}
}

type bookmarkSqlReader struct{}

func (r *bookmarkSqlReader) SelectByUserID(dh DataHandler, result *[]int, userID int) error {
	q := "SELECT podcast_id FROM bookmarks WHERE user_id=$1 ORDER BY id DESC"
	return handleError(sqlx.Select(dh, result, q, userID), q)
}

type bookmarkSqlWriter struct{}

func (db *bookmarkSqlWriter) Create(dh DataHandler, podcastID, userID int) error {
	q := "INSERT INTO bookmarks(podcast_id, user_id) VALUES($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return handleError(err, q)
}

func (db *bookmarkSqlWriter) Delete(dh DataHandler, podcastID, userID int) error {
	q := "DELETE FROM bookmarks WHERE podcast_id=$1 AND user_id=$2"
	_, err := dh.Exec(q, podcastID, userID)
	return handleError(err, q)
}

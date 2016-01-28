package store

import (
	"github.com/jmoiron/sqlx"
)

type BookmarkWriter interface {
	Create(DataHandler, int64, int64) error
	Delete(DataHandler, int64, int64) error
}

type BookmarkReader interface {
	SelectByUserID(DataHandler, *[]int64, int64) error
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

func (r *bookmarkSqlReader) SelectByUserID(dh DataHandler, result *[]int64, userID int64) error {
	q := "SELECT podcast_id FROM bookmarks WHERE user_id=$1"
	return sqlx.Select(dh, result, q, userID)
}

type bookmarkSqlWriter struct{}

func (db *bookmarkSqlWriter) Create(dh DataHandler, podcastID, userID int64) error {
	q := "INSERT INTO bookmarks(podcast_id, user_id) VALUES($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return err
}

func (db *bookmarkSqlWriter) Delete(dh DataHandler, podcastID, userID int64) error {
	q := "DELETE FROM bookmarks WHERE podcast_id=$1 AND user_id=$2"
	_, err := dh.Exec(q, podcastID, userID)
	return err
}

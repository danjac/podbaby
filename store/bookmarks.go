package store

import (
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

// BookmarkWriter handles all writes to bookmark data store
type BookmarkWriter interface {
	Create(DataHandler, int, int) error
	Delete(DataHandler, int, int) error
}

// BookmarkReader handles all reads from bookmark data store
type BookmarkReader interface {
	SelectByUserID(DataHandler, *[]int, int) error
}

// BookmarkStore handles manages interactions with bookmark data store
type BookmarkStore interface {
	BookmarkReader
	BookmarkWriter
}

type bookmarkSQLStore struct {
	BookmarkReader
	BookmarkWriter
}

func newBookmarkStore() BookmarkStore {
	return &bookmarkSQLStore{
		BookmarkReader: &bookmarkSQLReader{},
		BookmarkWriter: &bookmarkSQLWriter{},
	}
}

type bookmarkSQLReader struct{}

func (r *bookmarkSQLReader) SelectByUserID(dh DataHandler, result *[]int, userID int) error {
	q := "SELECT podcast_id FROM bookmarks WHERE user_id=$1 ORDER BY id DESC"
	return handleError(sqlx.Select(dh, result, q, userID), q)
}

type bookmarkSQLWriter struct{}

func (db *bookmarkSQLWriter) Create(dh DataHandler, podcastID, userID int) error {
	q := "INSERT INTO bookmarks(podcast_id, user_id) VALUES($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return handleError(err, q)
}

func (db *bookmarkSQLWriter) Delete(dh DataHandler, podcastID, userID int) error {
	q := "DELETE FROM bookmarks WHERE podcast_id=$1 AND user_id=$2"
	_, err := dh.Exec(q, podcastID, userID)
	return handleError(err, q)
}

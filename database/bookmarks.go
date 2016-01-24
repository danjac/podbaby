package database

import (
	"github.com/jmoiron/sqlx"
)

type BookmarkWriter interface {
	Create(DataHandler, int64, int64) error
	Delete(DataHandler, int64, int64) error
}

type BookmarkReader interface {
	SelectByUserID(DataHandler, int64) ([]int64, error)
}

type BookmarkDB interface {
	BookmarkReader
	BookmarkWriter
}

type BookmarkSqlDB struct {
	BookmarkReader
	BookmarkWriter
}

func newBookmarkDB() BookmarkDB {
	return &BookmarkSqlDB{
		BookmarkReader: &BookmarkSqlReader{},
		BookmarkWriter: &BookmarkSqlWriter{},
	}
}

type BookmarkSqlReader struct{}

func (r *BookmarkSqlReader) SelectByUserID(dh DataHandler, userID int64) ([]int64, error) {
	q := "SELECT podcast_id FROM bookmarks WHERE user_id=$1"
	var result []int64
	err := sqlx.Select(dh, &result, q, userID)
	return result, dbErr(err, q)
}

type BookmarkSqlWriter struct{}

func (db *BookmarkSqlWriter) Create(dh DataHandler, podcastID, userID int64) error {
	q := "INSERT INTO bookmarks(podcast_id, user_id) VALUES($1, $2)"
	_, err := dh.Exec(q, podcastID, userID)
	return dbErr(err, q)
}

func (db *BookmarkSqlWriter) Delete(dh DataHandler, podcastID, userID int64) error {
	q := "DELETE FROM bookmarks WHERE podcast_id=$1 AND user_id=$2"
	_, err := dh.Exec(q, podcastID, userID)
	return dbErr(err, q)
}

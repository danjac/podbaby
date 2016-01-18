package database

import (
	"github.com/jmoiron/sqlx"
)

type BookmarkWriter interface {
	Create(int64, int64) error
	Delete(int64, int64) error
}

type BookmarkReader interface {
	SelectByUserID(int64) ([]int64, error)
}

type BookmarkDB struct {
	BookmarkReader
	BookmarkWriter
}

func newBookmarkDB(db sqlx.Ext) *BookmarkDB {
	return &BookmarkDB{
		BookmarkReader: &BookmarkDBReader{db},
		BookmarkWriter: &BookmarkDBWriter{db},
	}
}

type BookmarkDBReader struct {
	sqlx.Ext
}

func (db *BookmarkDBReader) SelectByUserID(userID int64) ([]int64, error) {
	q := "SELECT podcast_id FROM bookmarks WHERE user_id=$1"
	var result []int64
	err := sqlx.Select(db, &result, q, userID)
	return result, dbErr(err, q)
}

type BookmarkDBWriter struct {
	sqlx.Ext
}

func (db *BookmarkDBWriter) Create(podcastID, userID int64) error {
	q := "INSERT INTO bookmarks(podcast_id, user_id) VALUES($1, $2)"
	_, err := db.Exec(q, podcastID, userID)
	return dbErr(err, q)
}

func (db *BookmarkDBWriter) Delete(podcastID, userID int64) error {
	q := "DELETE FROM bookmarks WHERE podcast_id=$1 AND user_id=$2"
	_, err := db.Exec(q, podcastID, userID)
	return dbErr(err, q)
}

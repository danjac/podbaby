package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
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

func newBookmarkDB(db sqlx.Ext, ps purse.Purse) *BookmarkDB {
	return &BookmarkDB{
		BookmarkReader: &BookmarkDBReader{db, ps},
		BookmarkWriter: &BookmarkDBWriter{db, ps},
	}
}

type BookmarkDBReader struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *BookmarkDBReader) SelectByUserID(userID int64) ([]int64, error) {
	q, _ := db.ps.Get("select_bookmarks.sql")
	var result []int64
	err := sqlx.Select(db, &result, q, userID)
	return result, sqlErr(err, q)
}

type BookmarkDBWriter struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *BookmarkDBWriter) Create(podcastID, userID int64) error {
	q, _ := db.ps.Get("insert_bookmark.sql")
	_, err := db.Exec(q, podcastID, userID)
	return sqlErr(err, q)
}

func (db *BookmarkDBWriter) Delete(podcastID, userID int64) error {
	q, _ := db.ps.Get("delete_bookmark.sql")
	_, err := db.Exec(q, podcastID, userID)
	return sqlErr(err, q)
}

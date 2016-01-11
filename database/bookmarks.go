package database

import (
	"github.com/danjac/podbaby/sql"
	"github.com/jmoiron/sqlx"
)

type BookmarkDB interface {
	Create(int64, int64) error
	Delete(int64, int64) error
	SelectByUserID(int64) ([]int64, error)
}

type defaultBookmarkDBImpl struct {
	*sqlx.DB
}

func (db *defaultBookmarkDBImpl) SelectByUserID(userID int64) ([]int64, error) {
	q, _ := sql.Queries.Get("select_bookmarks.sql")
	var result []int64
	err := db.Select(&result, q, userID)
	return result, err
}

func (db *defaultBookmarkDBImpl) Create(podcastID, userID int64) error {
	q, _ := sql.Queries.Get("insert_bookmark.sql")
	_, err := db.Exec(q, podcastID, userID)
	return err
}

func (db *defaultBookmarkDBImpl) Delete(podcastID, userID int64) error {
	q, _ := sql.Queries.Get("delete_bookmark.sql")
	_, err := db.Exec(q, podcastID, userID)
	return err
}

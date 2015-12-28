package database

import (
	"github.com/jmoiron/sqlx"
)

type BookmarkDB interface {
	Create(int64, int64) error
	Delete(int64, int64) error
}

type defaultBookmarkDBImpl struct {
	*sqlx.DB
}

func (db *defaultBookmarkDBImpl) Create(podcastID, userID int64) error {
	sql := "INSERT INTO bookmarks(podcast_id, user_id) VALUES($1, $2)"
	_, err := db.Exec(sql, podcastID, userID)
	return err
}

func (db *defaultBookmarkDBImpl) Delete(podcastID, userID int64) error {
	sql := "DELETE FROM bookmarks WHERE podcast_id=$1 AND user_id=$2"
	_, err := db.Exec(sql, podcastID, userID)
	return err
}

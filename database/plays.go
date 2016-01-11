package database

import (
	"github.com/danjac/podbaby/sql"
	"github.com/jmoiron/sqlx"
)

type PlayDB interface {
	Create(int64, int64) error
	DeleteAll(int64) error
}

type defaultPlayDBImpl struct {
	*sqlx.DB
}

func (db *defaultPlayDBImpl) Create(podcastID, userID int64) error {
	q, _ := sql.Queries.Get("add_play.sql")
	_, err := db.Exec(q, podcastID, userID)
	return err
}

func (db *defaultPlayDBImpl) DeleteAll(userID int64) error {
	q, _ := sql.Queries.Get("delete_plays_by_user_id.sql")
	_, err := db.Exec(q, userID)
	return err
}

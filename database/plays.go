package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
)

type PlayDB interface {
	Create(int64, int64) error
	DeleteAll(int64) error
}

type defaultPlayDBImpl struct {
	*sqlx.DB
	ps purse.Purse
}

func (db *defaultPlayDBImpl) Create(podcastID, userID int64) error {
	q, _ := db.ps.Get("add_play.sql")
	_, err := db.Exec(q, podcastID, userID)
	return sqlErr(err, q)
}

func (db *defaultPlayDBImpl) DeleteAll(userID int64) error {
	q, _ := db.ps.Get("delete_plays_by_user_id.sql")
	_, err := db.Exec(q, userID)
	return sqlErr(err, q)
}

package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
)

type PlayWriter interface {
	Create(int64, int64) error
	DeleteAll(int64) error
}

type PlayDB struct {
	PlayWriter
}

func newPlayDB(db sqlx.Ext, ps purse.Purse) *PlayDB {
	return &PlayDB{
		PlayWriter: &PlayDBWriter{db, ps},
	}
}

type PlayDBWriter struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *PlayDBWriter) Create(podcastID, userID int64) error {
	q, _ := db.ps.Get("add_play.sql")
	_, err := db.Exec(q, podcastID, userID)
	return sqlErr(err, q)
}

func (db *PlayDBWriter) DeleteAll(userID int64) error {
	q, _ := db.ps.Get("delete_plays_by_user_id.sql")
	_, err := db.Exec(q, userID)
	return sqlErr(err, q)
}

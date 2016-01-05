package database

import (
	"github.com/jmoiron/sqlx"
)

type PlayDB interface {
	Create(int64, int64) error
}

type defaultPlayDBImpl struct {
	*sqlx.DB
}

func (db *defaultPlayDBImpl) Create(podcastID, userID int64) error {
	_, err := db.Exec("SELECT add_play($1, $2)", podcastID, userID)
	return err
}

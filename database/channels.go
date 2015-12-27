package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// ChannelDB database model
type ChannelDB interface {
	GetOrCreate(*models.Channel) (bool, error)
}

type defaultChannelDBImpl struct {
	*sqlx.DB
}

func (db *defaultChannelDBImpl) GetOrCreate(ch *models.Channel) (bool, error) {

	if err := db.Get(channel, "SELECT * FROM channels WHERE url=$1", channel.URL); err == nil {
		return false, nil
	}

	query, args, err := sqlx.Named(`
    INSERT INTO channels (url, title, image, description)
    VALUES (:url, :title, :image, :description)
    RETURNING id`, ch)

	if err != nil {
		return false, err
	}

	return true, db.QueryRow(db.Rebind(query), args...).Scan(&ch.ID)
}

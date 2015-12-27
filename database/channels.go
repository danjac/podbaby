package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// ChannelDB database model
type ChannelDB interface {
	Create(*models.Channel) error
}

type defaultChannelDBImpl struct {
	*sqlx.DB
}

func (db *defaultChannelDBImpl) Create(ch *models.Channel) error {

	query, args, err := sqlx.Named("SELECT upsert_channel (:url, :title, :description, :image)", ch)

	if err != nil {
		return err
	}

	return db.QueryRow(db.Rebind(query), args...).Scan(&ch.ID)
}

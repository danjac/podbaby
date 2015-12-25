package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type ChannelDB interface {
	Create(*models.Channel) error
}

type defaultChannelDBImpl struct {
	*sqlx.DB
}

func (db *defaultChannelDBImpl) Create(ch *models.Channel) error {

	query, args, err := sqlx.Named(`
    INSERT INTO channels (url, title, image, description)  
    VALUES (:url, :title, :image, :description)
    RETURNING id`, ch)

	if err != nil {
		return err
	}

	return db.QueryRow(db.Rebind(query), args...).Scan(&ch.ID)
}

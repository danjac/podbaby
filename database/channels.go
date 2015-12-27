package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// ChannelDB database model
type ChannelDB interface {
	SelectSubscribed(int64) ([]models.Channel, error)
	GetByID(int64, int64) (*models.Channel, error)
	Create(*models.Channel) error
}

type defaultChannelDBImpl struct {
	*sqlx.DB
}

func (db *defaultChannelDBImpl) SelectSubscribed(userID int64) ([]models.Channel, error) {
	sql := `SELECT DISTINCT c.* FROM channels c
	JOIN subscriptions s ON s.channel_id = c.id
	WHERE s.user_id=$1 AND title IS NOT NULL
	ORDER BY title`
	var channels []models.Channel
	return channels, db.Select(&channels, sql, userID)
}

func (db *defaultChannelDBImpl) GetByID(id int64, userID int64) (*models.Channel, error) {
	sql := `SELECT channels.*,
        EXISTS(SELECT id FROM subscriptions WHERE channel_id=channels.id AND user_id=$1) AS is_subscribed
        FROM channels
        WHERE id=$2`
	channel := &models.Channel{}
	return channel, db.Get(channel, sql, userID, id)
}

func (db *defaultChannelDBImpl) Create(ch *models.Channel) error {

	query, args, err := sqlx.Named("SELECT upsert_channel (:url, :title, :description, :image)", ch)

	if err != nil {
		return err
	}

	return db.QueryRow(db.Rebind(query), args...).Scan(&ch.ID)
}

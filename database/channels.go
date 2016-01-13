package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
)

// ChannelDB database model
type ChannelDB interface {
	SelectAll() ([]models.Channel, error)
	SelectSubscribed(int64) ([]models.Channel, error)
	Search(string) ([]models.Channel, error)
	GetByID(int64) (*models.Channel, error)
	GetByURL(string) (*models.Channel, error)
	Create(*models.Channel) error
}

type defaultChannelDBImpl struct {
	*sqlx.DB
	ps purse.Purse
}

func (db *defaultChannelDBImpl) SelectAll() ([]models.Channel, error) {

	q, _ := db.ps.Get("select_all_channels.sql")
	var channels []models.Channel
	return channels, sqlErr(db.Select(&channels, q), q)
}

func (db *defaultChannelDBImpl) SelectSubscribed(userID int64) ([]models.Channel, error) {

	q, _ := db.ps.Get("select_subscribed_channels.sql")
	var channels []models.Channel
	return channels, sqlErr(db.Select(&channels, q, userID), q)
}

func (db *defaultChannelDBImpl) Search(query string) ([]models.Channel, error) {

	q, _ := db.ps.Get("search_channels.sql")

	var channels []models.Channel
	return channels, sqlErr(db.Select(&channels, q, query), q)
}

func (db *defaultChannelDBImpl) GetByURL(url string) (*models.Channel, error) {
	q, _ := db.ps.Get("get_channel_by_url.sql")
	channel := &models.Channel{}
	return channel, sqlErr(db.Get(channel, q, url), q)
}

func (db *defaultChannelDBImpl) GetByID(id int64) (*models.Channel, error) {
	q, _ := db.ps.Get("get_channel_by_id.sql")
	channel := &models.Channel{}
	return channel, sqlErr(db.Get(channel, q, id), q)
}

func (db *defaultChannelDBImpl) Create(ch *models.Channel) error {

	q, _ := db.ps.Get("upsert_channel.sql")

	q, args, err := sqlx.Named(q, ch)
	if err != nil {
		return sqlErr(err, q)
	}

	return sqlErr(db.QueryRowx(db.Rebind(q), args...).Scan(&ch.ID), q)
}

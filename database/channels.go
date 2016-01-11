package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/sql"
	"github.com/jmoiron/sqlx"
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
}

func (db *defaultChannelDBImpl) SelectAll() ([]models.Channel, error) {

	q, _ := sql.Queries.Get("select_all_channels.sql")
	var channels []models.Channel
	return channels, db.Select(&channels, q)
}

func (db *defaultChannelDBImpl) SelectSubscribed(userID int64) ([]models.Channel, error) {

	q, _ := sql.Queries.Get("select_subscribed_channels.sql")
	var channels []models.Channel
	return channels, db.Select(&channels, q, userID)
}

func (db *defaultChannelDBImpl) Search(query string) ([]models.Channel, error) {

	q, _ := sql.Queries.Get("search_channels.sql")

	var channels []models.Channel
	return channels, db.Select(&channels, q, query)
}

func (db *defaultChannelDBImpl) GetByURL(url string) (*models.Channel, error) {
	q, _ := sql.Queries.Get("get_channel_by_url.sql")
	channel := &models.Channel{}
	return channel, db.Get(channel, q, url)
}

func (db *defaultChannelDBImpl) GetByID(id int64) (*models.Channel, error) {
	q, _ := sql.Queries.Get("get_channel_by_id.sql")
	channel := &models.Channel{}
	return channel, db.Get(channel, q, id)
}

func (db *defaultChannelDBImpl) Create(ch *models.Channel) error {

	q, _ := sql.Queries.Get("upsert_channel.sql")

	q, args, err := sqlx.Named(q, ch)
	if err != nil {
		return err
	}

	return db.QueryRowx(db.Rebind(q), args...).Scan(&ch.ID)
}

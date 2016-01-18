package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
)

type ChannelReader interface {
	SelectAll() ([]models.Channel, error)
	SelectSubscribed(int64) ([]models.Channel, error)
	Search(string) ([]models.Channel, error)
	GetByID(int64) (*models.Channel, error)
	GetByURL(string) (*models.Channel, error)
}

type ChannelWriter interface {
	Create(*models.Channel) error
}

type ChannelDB struct {
	ChannelReader
	ChannelWriter
}

func newChannelDB(db sqlx.Ext, ps purse.Purse) *ChannelDB {
	return &ChannelDB{
		ChannelReader: &ChannelDBReader{db, ps},
		ChannelWriter: &ChannelDBWriter{db, ps},
	}
}

type ChannelDBReader struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *ChannelDBReader) SelectAll() ([]models.Channel, error) {
	q, _ := db.ps.Get("select_all_channels.sql")
	var channels []models.Channel
	return channels, sqlErr(sqlx.Select(db, &channels, q), q)
}

func (db *ChannelDBReader) SelectSubscribed(userID int64) ([]models.Channel, error) {

	q, _ := db.ps.Get("select_subscribed_channels.sql")
	var channels []models.Channel
	return channels, sqlErr(sqlx.Select(db, &channels, q, userID), q)
}

func (db *ChannelDBReader) Search(query string) ([]models.Channel, error) {

	q, _ := db.ps.Get("search_channels.sql")

	var channels []models.Channel
	return channels, sqlErr(sqlx.Select(db, &channels, q, query), q)
}

func (db *ChannelDBReader) GetByURL(url string) (*models.Channel, error) {
	q, _ := db.ps.Get("get_channel_by_url.sql")
	channel := &models.Channel{}
	return channel, sqlErr(sqlx.Get(db, channel, q, url), q)
}

func (db *ChannelDBReader) GetByID(id int64) (*models.Channel, error) {
	q, _ := db.ps.Get("get_channel_by_id.sql")
	channel := &models.Channel{}
	return channel, sqlErr(sqlx.Get(db, channel, q, id), q)
}

type ChannelDBWriter struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *ChannelDBWriter) Create(ch *models.Channel) error {

	q, _ := db.ps.Get("upsert_channel.sql")

	q, args, err := sqlx.Named(q, ch)
	if err != nil {
		return sqlErr(err, q)
	}

	return sqlErr(db.QueryRowx(db.Rebind(q), args...).Scan(&ch.ID), q)
}

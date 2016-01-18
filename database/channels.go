package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
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

func newChannelDB(db sqlx.Ext) *ChannelDB {
	return &ChannelDB{
		ChannelReader: &ChannelDBReader{db},
		ChannelWriter: &ChannelDBWriter{db},
	}
}

type ChannelDBReader struct {
	sqlx.Ext
}

func (db *ChannelDBReader) SelectAll() ([]models.Channel, error) {
	q := `SELECT id, title, description, categories, url, image, website 
FROM channels`
	var channels []models.Channel
	return channels, sqlErr(sqlx.Select(db, &channels, q), q)
}

func (db *ChannelDBReader) SelectSubscribed(userID int64) ([]models.Channel, error) {

	q := `SELECT c.id, c.title, c.description, c.image, c.url, c.website
FROM channels c
JOIN subscriptions s ON s.channel_id = c.id
WHERE s.user_id=$1 AND title IS NOT NULL AND title != ''
GROUP BY c.id
ORDER BY title
`
	var channels []models.Channel
	return channels, sqlErr(sqlx.Select(db, &channels, q, userID), q)
}

func (db *ChannelDBReader) Search(query string) ([]models.Channel, error) {

	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c, plainto_tsquery($1) as q
WHERE (c.tsv @@ q)
ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 20`
	var channels []models.Channel
	return channels, sqlErr(sqlx.Select(db, &channels, q, query), q)
}

func (db *ChannelDBReader) GetByURL(url string) (*models.Channel, error) {
	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c
WHERE url=$1`
	channel := &models.Channel{}
	return channel, sqlErr(sqlx.Get(db, channel, q, url), q)
}

func (db *ChannelDBReader) GetByID(id int64) (*models.Channel, error) {
	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c
WHERE id=$1`
	channel := &models.Channel{}
	return channel, sqlErr(sqlx.Get(db, channel, q, id), q)
}

type ChannelDBWriter struct {
	sqlx.Ext
}

func (db *ChannelDBWriter) Create(ch *models.Channel) error {

	q := `SELECT upsert_channel (
:url, 
:title, 
:description, 
:image, 
:categories, 
:website
)`

	q, args, err := sqlx.Named(q, ch)
	if err != nil {
		return sqlErr(err, q)
	}

	return sqlErr(db.QueryRowx(db.Rebind(q), args...).Scan(&ch.ID), q)
}

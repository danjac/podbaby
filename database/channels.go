package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

const maxRecommendations = 20

type ChannelReader interface {
	SelectAll() ([]models.Channel, error)
	SelectByCategoryID(int64) ([]models.Channel, error)
	SelectSubscribed(int64) ([]models.Channel, error)
	SelectRelated(int64) ([]models.Channel, error)
	SelectRecommended() ([]models.Channel, error)
	SelectRecommendedByUserID(int64) ([]models.Channel, error)
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
	q := `SELECT id, title, description, url, image, website 
FROM channels`
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(db, &channels, q), q)
}

func (db *ChannelDBReader) SelectByCategoryID(categoryID int64) ([]models.Channel, error) {
	q := `SELECT c.id, c.title, c.image, c.description, c.website, c.url
FROM channels c
JOIN channels_categories cc 
ON cc.channel_id = c.id
WHERE cc.category_id=$1
GROUP BY c.id
ORDER BY c.title`

	var channels []models.Channel
	return channels, dbErr(sqlx.Select(db, &channels, q, categoryID), q)
}

func (db *ChannelDBReader) SelectRelated(channelID int64) ([]models.Channel, error) {
	q := `SELECT c.id, c.title, c.image, c.description, c.website, c.url
FROM channels c
JOIN subscriptions s ON s.channel_id=c.id
WHERE s.user_id in (
  SELECT user_id FROM subscriptions WHERE channel_id=$1
) AND s.channel_id != $1
GROUP BY c.id
ORDER BY RANDOM() DESC LIMIT 3`

	var channels []models.Channel
	return channels, dbErr(sqlx.Select(db, &channels, q, channelID), q)
}

func (db *ChannelDBReader) SelectRecommended() ([]models.Channel, error) {
	q := `SELECT c.id, c.title, c.image, c.description, c.website, c.url
FROM channels c
JOIN subscriptions s ON s.channel_id = c.id
GROUP BY c.id
ORDER BY COUNT(DISTINCT(s.id)) DESC LIMIT $1
    `
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(db, &channels, q, maxRecommendations), q)
}

func (db *ChannelDBReader) SelectRecommendedByUserID(userID int64) ([]models.Channel, error) {
	q := `
WITH user_subs AS (SELECT channel_id FROM subscriptions WHERE user_id=$1)
SELECT c.id, c.title, c.description, c.image, c.url, c.website
FROM channels c
JOIN channels_categories cc ON cc.channel_id=c.id
WHERE cc.category_id IN (
   SELECT cc.category_id FROM channels_categories cc
   WHERE cc.channel_id IN (SELECT channel_id FROM user_subs)
)
AND c.id NOT IN (SELECT channel_id FROM user_subs)
GROUP BY c.id
ORDER BY RANDOM()
LIMIT $2`
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(db, &channels, q, userID, maxRecommendations), q)
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
	return channels, dbErr(sqlx.Select(db, &channels, q, userID), q)
}

func (db *ChannelDBReader) Search(query string) ([]models.Channel, error) {

	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c, plainto_tsquery($1) as q
WHERE (c.tsv @@ q)
ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 20`
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(db, &channels, q, query), q)
}

func (db *ChannelDBReader) GetByURL(url string) (*models.Channel, error) {
	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c
WHERE url=$1`
	channel := &models.Channel{}
	return channel, dbErr(sqlx.Get(db, channel, q, url), q)
}

func (db *ChannelDBReader) GetByID(id int64) (*models.Channel, error) {
	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c
WHERE id=$1`
	channel := &models.Channel{}
	return channel, dbErr(sqlx.Get(db, channel, q, id), q)
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
    :keywords, 
    :website
)`

	q, args, err := sqlx.Named(q, ch)
	if err != nil {
		return dbErr(err, q)
	}

	return dbErr(db.QueryRowx(db.Rebind(q), args...).Scan(&ch.ID), q)
}

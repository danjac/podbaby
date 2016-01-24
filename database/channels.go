package database

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

const maxRecommendations = 20

type ChannelReader interface {
	SelectAll(DataHandler) ([]models.Channel, error)
	SelectByCategoryID(DataHandler, int64) ([]models.Channel, error)
	SelectSubscribed(DataHandler, int64) ([]models.Channel, error)
	SelectRelated(DataHandler, int64) ([]models.Channel, error)
	SelectRecommended(DataHandler) ([]models.Channel, error)
	SelectRecommendedByUserID(DataHandler, int64) ([]models.Channel, error)
	Search(DataHandler, string) ([]models.Channel, error)
	GetByID(DataHandler, int64) (*models.Channel, error)
	GetByURL(DataHandler, string) (*models.Channel, error)
}

type ChannelWriter interface {
	Create(DataHandler, *models.Channel) error
	AddCategories(DataHandler, *models.Channel) error
	AddPodcasts(DataHandler, *models.Channel) error
}

type ChannelDB interface {
	ChannelReader
	ChannelWriter
}

type ChannelSqlDB struct {
	ChannelReader
	ChannelWriter
}

func newChannelDB() ChannelDB {
	return &ChannelSqlDB{
		ChannelReader: &ChannelSqlReader{},
		ChannelWriter: &ChannelSqlWriter{},
	}
}

type ChannelSqlReader struct{}

func (r *ChannelSqlReader) SelectAll(dh DataHandler) ([]models.Channel, error) {
	q := `SELECT id, title, description, url, image, website 
FROM channels`
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(dh, &channels, q), q)
}

func (r *ChannelSqlReader) SelectByCategoryID(dh DataHandler, categoryID int64) ([]models.Channel, error) {
	q := `SELECT c.id, c.title, c.image, c.description, c.website, c.url
FROM channels c
JOIN channels_categories cc 
ON cc.channel_id = c.id
WHERE cc.category_id=$1
GROUP BY c.id
ORDER BY c.title`

	var channels []models.Channel
	return channels, dbErr(sqlx.Select(dh, &channels, q, categoryID), q)
}

func (r *ChannelSqlReader) SelectRelated(dh DataHandler, channelID int64) ([]models.Channel, error) {
	q := `SELECT c.id, c.title, c.image, c.description, c.website, c.url
FROM channels c
JOIN subscriptions s ON s.channel_id=c.id
WHERE s.user_id in (
  SELECT user_id FROM subscriptions WHERE channel_id=$1
) AND s.channel_id != $1
GROUP BY c.id
ORDER BY RANDOM() DESC LIMIT 3`

	var channels []models.Channel
	return channels, dbErr(sqlx.Select(dh, &channels, q, channelID), q)
}

func (r *ChannelSqlReader) SelectRecommended(dh DataHandler) ([]models.Channel, error) {
	q := `SELECT c.id, c.title, c.image, c.description, c.website, c.url
FROM channels c
JOIN subscriptions s ON s.channel_id = c.id
GROUP BY c.id
ORDER BY COUNT(DISTINCT(s.id)) DESC LIMIT $1
    `
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(dh, &channels, q, maxRecommendations), q)
}

func (r *ChannelSqlReader) SelectRecommendedByUserID(dh DataHandler, userID int64) ([]models.Channel, error) {
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
	return channels, dbErr(sqlx.Select(dh, &channels, q, userID, maxRecommendations), q)
}

func (r *ChannelSqlReader) SelectSubscribed(dh DataHandler, userID int64) ([]models.Channel, error) {

	q := `SELECT c.id, c.title, c.description, c.image, c.url, c.website
FROM channels c
JOIN subscriptions s ON s.channel_id = c.id
WHERE s.user_id=$1 AND title IS NOT NULL AND title != ''
GROUP BY c.id
ORDER BY title
`
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(dh, &channels, q, userID), q)
}

func (r *ChannelSqlReader) Search(dh DataHandler, query string) ([]models.Channel, error) {

	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c, plainto_tsquery($1) as q
WHERE (c.tsv @@ q)
ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 20`
	var channels []models.Channel
	return channels, dbErr(sqlx.Select(dh, &channels, q, query), q)
}

func (r *ChannelSqlReader) GetByURL(dh DataHandler, url string) (*models.Channel, error) {
	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c
WHERE url=$1`
	channel := &models.Channel{}
	return channel, dbErr(sqlx.Get(dh, channel, q, url), q)
}

func (r *ChannelSqlReader) GetByID(dh DataHandler, id int64) (*models.Channel, error) {
	q := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c
WHERE id=$1`
	channel := &models.Channel{}
	return channel, dbErr(sqlx.Get(dh, channel, q, id), q)
}

type ChannelSqlWriter struct{}

func (w *ChannelSqlWriter) Create(dh DataHandler, ch *models.Channel) error {

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

	return dbErr(dh.QueryRowx(dh.Rebind(q), args...).Scan(&ch.ID), q)
}

func (w *ChannelSqlWriter) AddCategories(dh DataHandler, channel *models.Channel) error {
	if len(channel.Categories) == 0 {
		return nil
	}
	args := []interface{}{
		channel.ID,
	}

	params := make([]string, 0, len(channel.Categories))
	for i, category := range channel.Categories {
		params = append(params, fmt.Sprintf("$%v", i+2))
		args = append(args, category)
	}

	q := fmt.Sprintf("SELECT add_categories($1, ARRAY[%s])", strings.Join(params, ", "))
	_, err := dh.Exec(q, args...)
	return dbErr(err, q)
}

func (w *ChannelSqlWriter) AddPodcasts(dh DataHandler, channel *models.Channel) error {

	q := `SELECT insert_podcast(
        :channel_id, 
        :guid,
        :title, 
        :description, 
        :enclosure_url, 
        :source,
        :pub_date)`

	stmt, err := dh.PrepareNamed(dh.Rebind(q))
	if err != nil {
		return err
	}

	for _, pc := range channel.Podcasts {
		pc.ChannelID = channel.ID
		err = dbErr(stmt.QueryRowx(&pc).Scan(&pc.ID), q)
		if err != nil {
			return err
		}
	}
	return nil

}

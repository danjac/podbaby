package store

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

const maxRecommendations = 20

type ChannelReader interface {
	SelectAll(DataHandler, *[]models.Channel) error
	SelectByCategoryID(DataHandler, *[]models.Channel, int64) error
	SelectSubscribed(DataHandler, *[]models.Channel, int64) error
	SelectRelated(DataHandler, *[]models.Channel, int64) error
	SelectRecommended(DataHandler, *[]models.Channel) error
	SelectRecommendedByUserID(DataHandler, *[]models.Channel, int64) error
	Search(DataHandler, *[]models.Channel, string) error
	GetByID(DataHandler, *models.Channel, int64) error
	GetByURL(DataHandler, *models.Channel, string) error
}

type ChannelWriter interface {
	Create(DataHandler, *models.Channel) error
	AddCategories(DataHandler, *models.Channel) error
	AddPodcasts(DataHandler, *models.Channel) error
}

type ChannelStore interface {
	ChannelReader
	ChannelWriter
}

type channelSqlStore struct {
	ChannelReader
	ChannelWriter
}

func newChannelStore() ChannelStore {
	return &channelSqlStore{
		ChannelReader: &channelSqlReader{},
		ChannelWriter: &channelSqlWriter{},
	}
}

type channelSqlReader struct{}

func (r *channelSqlReader) SelectAll(dh DataHandler, channels *[]models.Channel) error {
	q := "SELECT id, title, description, url, image, website, num_podcasts FROM channels"
	return sqlx.Select(dh, &channels, q)
}

func (r *channelSqlReader) SelectByCategoryID(dh DataHandler, channels *[]models.Channel, categoryID int64) error {
	q := `
    SELECT c.id, c.title, c.image, c.description, c.website, c.url, c.num_podcasts
    FROM channels c
    JOIN channels_categories cc 
    ON cc.channel_id = c.id
    WHERE cc.category_id=$1
    GROUP BY c.id
    ORDER BY c.title`
	return sqlx.Select(dh, channels, q, categoryID)
}

func (r *channelSqlReader) SelectRelated(dh DataHandler, channels *[]models.Channel, channelID int64) error {
	q := `
    SELECT c.id, c.title, c.image, c.description, c.website, c.url, c.num_podcasts
    FROM channels c
    JOIN subscriptions s ON s.channel_id=c.id
    WHERE s.user_id in (
      SELECT user_id FROM subscriptions WHERE channel_id=$1
    ) AND s.channel_id != $1
    GROUP BY c.id
    ORDER BY RANDOM() DESC LIMIT 3`

	return sqlx.Select(dh, channels, q, channelID)
}

func (r *channelSqlReader) SelectRecommended(dh DataHandler, channels *[]models.Channel) error {
	q := `
    SELECT c.id, c.title, c.image, c.description, c.website, c.url, c.num_podcasts
    FROM channels c
    JOIN subscriptions s ON s.channel_id = c.id
    GROUP BY c.id
    ORDER BY COUNT(DISTINCT(s.id)) DESC LIMIT $1
    `
	return sqlx.Select(dh, channels, q, maxRecommendations)
}

func (r *channelSqlReader) SelectRecommendedByUserID(dh DataHandler, channels *[]models.Channel, userID int64) error {
	q := `
    WITH user_subs AS (SELECT channel_id FROM subscriptions WHERE user_id=$1)
    SELECT c.id, c.title, c.description, c.image, c.url, c.website, c.num_podcasts
    FROM channels c
    JOIN channels_categories cc ON cc.channel_id=c.id
    WHERE (cc.category_id IN (
       SELECT cc.category_id FROM channels_categories cc
       WHERE cc.channel_id IN (SELECT channel_id FROM user_subs)
    ) AND c.id NOT IN (SELECT channel_id FROM user_subs)
    ) OR ((SELECT COUNT(channel_id) FROM user_subs) = 0)
    GROUP BY c.id
    ORDER BY RANDOM()
    LIMIT $2`
	return sqlx.Select(dh, channels, q, userID, maxRecommendations)
}

func (r *channelSqlReader) SelectSubscribed(dh DataHandler, channels *[]models.Channel, userID int64) error {

	q := `
    SELECT c.id, c.title, c.description, c.image, c.url, c.website, c.num_podcasts
    FROM channels c
    JOIN subscriptions s ON s.channel_id = c.id
    WHERE s.user_id=$1 AND title IS NOT NULL AND title != ''
    GROUP BY c.id
    ORDER BY title`
	return sqlx.Select(dh, channels, q, userID)
}

func (r *channelSqlReader) Search(dh DataHandler, channels *[]models.Channel, query string) error {

	q := `
    SELECT c.id, c.title, c.description, c.url, c.image, c.website, c.num_podcasts
    FROM channels c, plainto_tsquery($1) as q
    WHERE (c.tsv @@ q)
    ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 20`
	return sqlx.Select(dh, channels, q, query)
}

func (r *channelSqlReader) GetByURL(dh DataHandler, channel *models.Channel, url string) error {
	q := `
    SELECT id, title, description, url, image, website, num_podcasts
    FROM channels
    WHERE url=$1`
	return sqlx.Get(dh, channel, q, url)
}

func (r *channelSqlReader) GetByID(dh DataHandler, channel *models.Channel, id int64) error {
	q := `
    SELECT c.id, c.title, c.description, c.url, c.image, c.website, c.num_podcasts
    FROM channels c
    WHERE id=$1`
	return sqlx.Get(dh, channel, q, id)
}

type channelSqlWriter struct{}

func (w *channelSqlWriter) Create(dh DataHandler, ch *models.Channel) error {

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
		return err
	}

	return dh.QueryRowx(dh.Rebind(q), args...).Scan(&ch.ID)
}

func (w *channelSqlWriter) AddCategories(dh DataHandler, channel *models.Channel) error {
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
	return err
}

func (w *channelSqlWriter) AddPodcasts(dh DataHandler, channel *models.Channel) error {

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
		err = stmt.QueryRowx(&pc).Scan(&pc.ID)
		if err != nil {
			return err
		}
	}
	return nil

}

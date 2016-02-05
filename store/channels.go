package store

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"strings"
)

const maxRecommendations = 20

// ChannelReader handles reads from channel data store
type ChannelReader interface {
	SelectAll(DataHandler, *[]models.Channel) error
	SelectByCategoryID(DataHandler, *[]models.Channel, int) error
	SelectSubscribed(DataHandler, *[]models.Channel, int) error
	SelectRelated(DataHandler, *[]models.Channel, int) error
	SelectRecommended(DataHandler, *[]models.Channel) error
	SelectRecommendedByUserID(DataHandler, *[]models.Channel, int) error
	Search(DataHandler, *[]models.Channel, string) error
	GetByID(DataHandler, *models.Channel, int) error
	GetByURL(DataHandler, *models.Channel, string) error
}

// ChannelWriter handles writes to channel data store
type ChannelWriter interface {
	// CreateOrUpdate inserts/updates channel, and adds new categories/channels if needed
	CreateOrUpdate(DataHandler, *models.Channel) error
	AddCategories(DataHandler, *models.Channel) error
	AddPodcasts(DataHandler, *models.Channel) error
}

// ChannelStore handles interactions with channel data store
type ChannelStore interface {
	ChannelReader
	ChannelWriter
}

type channelSQLStore struct {
	ChannelReader
	ChannelWriter
}

func newChannelStore() ChannelStore {
	return &channelSQLStore{
		ChannelReader: &channelSQLReader{},
		ChannelWriter: &channelSQLWriter{},
	}
}

type channelSQLReader struct{}

func (r *channelSQLReader) SelectAll(dh DataHandler, channels *[]models.Channel) error {
	q := "SELECT id, title, description, url, image, website, num_podcasts FROM channels"
	return handleError(sqlx.Select(dh, channels, q), q)
}

func (r *channelSQLReader) SelectByCategoryID(dh DataHandler, channels *[]models.Channel, categoryID int) error {
	q := `
    SELECT c.id, c.title, c.image, c.description, c.website, c.url, c.num_podcasts
    FROM channels c
    JOIN channels_categories cc 
    ON cc.channel_id = c.id
    WHERE cc.category_id=$1
    GROUP BY c.id
    ORDER BY c.title`
	return handleError(sqlx.Select(dh, channels, q, categoryID), q)
}

func (r *channelSQLReader) SelectRelated(dh DataHandler, channels *[]models.Channel, channelID int) error {
	q := `
    SELECT c.id, c.title, c.image, c.description, c.website, c.url, c.num_podcasts
    FROM channels c
    JOIN subscriptions s ON s.channel_id=c.id
    WHERE s.user_id in (
      SELECT user_id FROM subscriptions WHERE channel_id=$1
    ) AND s.channel_id != $1
    GROUP BY c.id
    ORDER BY RANDOM() DESC LIMIT 3`

	return handleError(sqlx.Select(dh, channels, q, channelID), q)
}

func (r *channelSQLReader) SelectRecommended(dh DataHandler, channels *[]models.Channel) error {
	q := `
    SELECT c.id, c.title, c.image, c.description, c.website, c.url, c.num_podcasts
    FROM channels c
    JOIN subscriptions s ON s.channel_id = c.id
    GROUP BY c.id
    ORDER BY COUNT(DISTINCT(s.id)) DESC LIMIT $1
    `
	return handleError(sqlx.Select(dh, channels, q, maxRecommendations), q)
}

func (r *channelSQLReader) SelectRecommendedByUserID(dh DataHandler, channels *[]models.Channel, userID int) error {
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
	return handleError(sqlx.Select(dh, channels, q, userID, maxRecommendations), q)
}

func (r *channelSQLReader) SelectSubscribed(dh DataHandler, channels *[]models.Channel, userID int) error {

	q := `
    SELECT c.id, c.title, c.description, c.image, c.url, c.website, c.num_podcasts
    FROM channels c
    JOIN subscriptions s ON s.channel_id = c.id
    WHERE s.user_id=$1 AND title IS NOT NULL AND title != ''
    GROUP BY c.id
    ORDER BY title`
	return handleError(sqlx.Select(dh, channels, q, userID), q)
}

func (r *channelSQLReader) Search(dh DataHandler, channels *[]models.Channel, query string) error {

	q := `
    SELECT c.id, c.title, c.description, c.url, c.image, c.website, c.num_podcasts
    FROM channels c, plainto_tsquery($1) as q
    WHERE (c.tsv @@ q)
    ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 20`
	return handleError(sqlx.Select(dh, channels, q, query), q)
}

func (r *channelSQLReader) GetByURL(dh DataHandler, channel *models.Channel, url string) error {
	q := `
    SELECT id, title, description, url, image, website, num_podcasts
    FROM channels
    WHERE url=$1`
	return handleError(sqlx.Get(dh, channel, q, url), q)
}

func (r *channelSQLReader) GetByID(dh DataHandler, channel *models.Channel, id int) error {
	q := `
    SELECT c.id, c.title, c.description, c.url, c.image, c.website, c.num_podcasts
    FROM channels c
    WHERE id=$1`
	return handleError(sqlx.Get(dh, channel, q, id), q)
}

type channelSQLWriter struct{}

func (w *channelSQLWriter) CreateOrUpdate(dh DataHandler, ch *models.Channel) error {

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
		return handleError(err, q)
	}

	if err := dh.QueryRowx(dh.Rebind(q), args...).Scan(&ch.ID); err != nil {
		return handleError(err, q)
	}

	if err = w.AddCategories(dh, ch); err != nil {
		return handleError(err, q)
	}

	if err = w.AddPodcasts(dh, ch); err != nil {
		return handleError(err, q)
	}

	return nil

}

func (w *channelSQLWriter) AddCategories(dh DataHandler, channel *models.Channel) error {
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

	q := fmt.Sprintf("SELECT add_categories ($1, ARRAY[%s])", strings.Join(params, ", "))
	_, err := dh.Exec(q, args...)
	return handleError(err, q)
}

func (w *channelSQLWriter) AddPodcasts(dh DataHandler, channel *models.Channel) error {

	q := `SELECT insert_podcast (
        :channel_id, 
        :guid,
        :title, 
        :description, 
        :enclosure_url, 
        :source,
        :pub_date)`

	stmt, err := dh.PrepareNamed(dh.Rebind(q))
	defer stmt.Close()

	if err != nil {
		return handleError(err, q)
	}

	for _, pc := range channel.Podcasts {
		pc.ChannelID = channel.ID
		err = stmt.QueryRowx(&pc).Scan(&pc.ID)
		if err != nil {
			return handleError(err, q)
		}
	}
	return nil

}

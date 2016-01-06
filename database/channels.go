package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// ChannelDB database model
type ChannelDB interface {
	SelectAll() ([]models.Channel, error)
	SelectSubscribed(int64) ([]models.Channel, error)
	Search(string, int64) ([]models.Channel, error)
	GetByID(int64, int64) (*models.Channel, error)
	GetByURL(string, int64) (*models.Channel, error)
	Create(*models.Channel) error
}

type defaultChannelDBImpl struct {
	*sqlx.DB
}

func (db *defaultChannelDBImpl) SelectAll() ([]models.Channel, error) {
	sql := "SELECT id, title, description, categories, url, image, website FROM channels"
	var channels []models.Channel
	return channels, db.Select(&channels, sql)
}

func (db *defaultChannelDBImpl) SelectSubscribed(userID int64) ([]models.Channel, error) {
	sql := `SELECT c.id, c.title, c.description, c.image, c.url, c.website,
    1 AS is_subscribed
    FROM channels c
    JOIN subscriptions s ON s.channel_id = c.id
    WHERE s.user_id=$1 AND title IS NOT NULL AND title != ''
    GROUP BY c.id
    ORDER BY title`
	var channels []models.Channel
	return channels, db.Select(&channels, sql, userID)
}

func (db *defaultChannelDBImpl) Search(query string, userID int64) ([]models.Channel, error) {

	sql := `SELECT c.id, c.title, c.description, c.url, c.image, c.website,
    EXISTS(SELECT id FROM subscriptions WHERE channel_id=c.id AND user_id=$1) AS is_subscribed
    FROM channels c, plainto_tsquery($2) as q
    WHERE (c.tsv @@ q)
    ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($2)) DESC LIMIT 10`

	var channels []models.Channel
	return channels, db.Select(&channels, sql, userID, query)
}

func (db *defaultChannelDBImpl) GetByURL(url string, userID int64) (*models.Channel, error) {
	sql := `SELECT c.id, c.title, c.description, c.url, c.image, c.website,
        EXISTS(SELECT id FROM subscriptions WHERE channel_id=c.id AND user_id=$1) AS is_subscribed
        FROM channels c
        WHERE url=$2`
	channel := &models.Channel{}
	return channel, db.Get(channel, sql, userID, url)
}

func (db *defaultChannelDBImpl) GetByID(id int64, userID int64) (*models.Channel, error) {
	sql := `SELECT c.id, c.title, c.description, c.url, c.image, c.website,
        EXISTS(SELECT id FROM subscriptions WHERE channel_id=c.id AND user_id=$1) AS is_subscribed
        FROM channels c
        WHERE id=$2`
	channel := &models.Channel{}
	return channel, db.Get(channel, sql, userID, id)
}

func (db *defaultChannelDBImpl) Create(ch *models.Channel) error {

	query, args, err := sqlx.Named(`
        SELECT upsert_channel (
            :url, 
            :title, 
            :description, 
            :image, 
            :categories, 
            :website)`, ch)

	if err != nil {
		return err
	}

	return db.QueryRowx(db.Rebind(query), args...).Scan(&ch.ID)
}

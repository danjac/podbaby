package database

import (
	"github.com/danjac/podbaby/models"
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
	sql := "SELECT id, title, description, categories, url, image, website FROM channels"
	var channels []models.Channel
	return channels, db.Select(&channels, sql)
}

func (db *defaultChannelDBImpl) SelectSubscribed(userID int64) ([]models.Channel, error) {
	sql := `SELECT c.id, c.title, c.description, c.image, c.url, c.website
    FROM channels c
    JOIN subscriptions s ON s.channel_id = c.id
    WHERE s.user_id=$1 AND title IS NOT NULL AND title != ''
    GROUP BY c.id
    ORDER BY title`
	var channels []models.Channel
	return channels, db.Select(&channels, sql, userID)
}

func (db *defaultChannelDBImpl) Search(query string) ([]models.Channel, error) {

	sql := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
    FROM channels c, plainto_tsquery($1) as q
    WHERE (c.tsv @@ q)
    ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 10`

	var channels []models.Channel
	return channels, db.Select(&channels, sql, query)
}

func (db *defaultChannelDBImpl) GetByURL(url string) (*models.Channel, error) {
	sql := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
        FROM channels c
        WHERE url=$1`
	channel := &models.Channel{}
	return channel, db.Get(channel, sql, url)
}

func (db *defaultChannelDBImpl) GetByID(id int64) (*models.Channel, error) {
	sql := `SELECT c.id, c.title, c.description, c.url, c.image, c.website
        FROM channels c
        WHERE id=$1`
	channel := &models.Channel{}
	return channel, db.Get(channel, sql, id)
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

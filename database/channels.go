package database

import (
	"fmt"
	"strings"

	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// ChannelDB database model
type ChannelDB interface {
	GetAll() ([]models.Channel, error)
	SelectSubscribed(int64) ([]models.Channel, error)
	Search(string, int64) ([]models.Channel, error)
	GetByID(int64, int64) (*models.Channel, error)
	Create(*models.Channel) error
}

type defaultChannelDBImpl struct {
	*sqlx.DB
}

func (db *defaultChannelDBImpl) GetAll() ([]models.Channel, error) {
	var channels []models.Channel
	return channels, db.Select(&channels, "SELECT * FROM channels")
}

func (db *defaultChannelDBImpl) SelectSubscribed(userID int64) ([]models.Channel, error) {
	sql := `SELECT DISTINCT c.* FROM channels c
  JOIN subscriptions s ON s.channel_id = c.id
  WHERE s.user_id=$1 AND title IS NOT NULL
  ORDER BY title`
	var channels []models.Channel
	return channels, db.Select(&channels, sql, userID)
}

func (db *defaultChannelDBImpl) Search(query string, userID int64) ([]models.Channel, error) {

	tokens := strings.Split(query, " ")

	var searchArgs []interface{}

	sql := `SELECT channels.*,
    EXISTS(SELECT id FROM subscriptions WHERE channel_id=channels.id AND user_id=$1) AS is_subscribed
    FROM channels WHERE`

	searchArgs = append(searchArgs, userID)

	var clauses []string

	for counter, token := range tokens {
		searchArgs = append(searchArgs, fmt.Sprintf("%%%s%%", token))
		clauses = append(clauses, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", counter+2, counter+2))
	}

	sql += " " + strings.Join(clauses, " AND ") + " ORDER BY title DESC LIMIT 20"

	var channels []models.Channel
	return channels, db.Select(&channels, sql, searchArgs...)
}

func (db *defaultChannelDBImpl) GetByID(id int64, userID int64) (*models.Channel, error) {
	sql := `SELECT channels.*,
        EXISTS(SELECT id FROM subscriptions WHERE channel_id=channels.id AND user_id=$1) AS is_subscribed
        FROM channels
        WHERE id=$2`
	channel := &models.Channel{}
	return channel, db.Get(channel, sql, userID, id)
}

func (db *defaultChannelDBImpl) Create(ch *models.Channel) error {

	query, args, err := sqlx.Named("SELECT upsert_channel (:url, :title, :description, :image)", ch)

	if err != nil {
		return err
	}

	return db.QueryRow(db.Rebind(query), args...).Scan(&ch.ID)
}

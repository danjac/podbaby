package database

import (
	"github.com/jmoiron/sqlx"
)

type SubscriptionDB interface {
	Create(int64, int64) error
	Delete(int64, int64) error
}

type defaultSubscriptionDBImpl struct {
	*sqlx.DB
}

func (db *defaultSubscriptionDBImpl) Create(channelID, userID int64) error {
	sql := "INSERT INTO subscriptions(channel_id, user_id) VALUES($1, $2)"
	_, err := db.Exec(sql, channelID, userID)
	return err
}

func (db *defaultSubscriptionDBImpl) Delete(channelID, userID int64) error {
	sql := "DELETE subscriptions WHERE channel_id=$1 AND user_id=$2 VALUES($1, $2)"
	_, err := db.Exec(sql, channelID, userID)
	return err
}

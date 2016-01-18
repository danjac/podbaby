package database

import (
	"github.com/jmoiron/sqlx"
)

type SubscriptionWriter interface {
	Create(int64, int64) error
	Delete(int64, int64) error
}

type SubscriptionReader interface {
	SelectByUserID(int64) ([]int64, error)
}

type SubscriptionDB struct {
	SubscriptionReader
	SubscriptionWriter
}

func newSubscriptionDB(db sqlx.Ext) *SubscriptionDB {
	return &SubscriptionDB{
		SubscriptionWriter: &SubscriptionDBWriter{db},
		SubscriptionReader: &SubscriptionDBReader{db},
	}
}

type SubscriptionDBReader struct {
	sqlx.Ext
}

func (db *SubscriptionDBReader) SelectByUserID(userID int64) ([]int64, error) {
	q := "SELECT channel_id FROM subscriptions WHERE user_id=$1"
	var result []int64
	err := sqlx.Select(db, &result, q, userID)
	return result, dbErr(err, q)
}

type SubscriptionDBWriter struct {
	sqlx.Ext
}

func (db *SubscriptionDBWriter) Create(channelID, userID int64) error {
	q := "INSERT INTO subscriptions(channel_id, user_id) VALUES($1, $2)"
	_, err := db.Exec(q, channelID, userID)
	return dbErr(err, q)
}

func (db *SubscriptionDBWriter) Delete(channelID, userID int64) error {
	q := "DELETE FROM subscriptions WHERE channel_id=$1 AND user_id=$2"
	_, err := db.Exec(q, channelID, userID)
	return dbErr(err, q)
}

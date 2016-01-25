package store

import (
	"github.com/jmoiron/sqlx"
)

type SubscriptionWriter interface {
	Create(DataHandler, int64, int64) error
	Delete(DataHandler, int64, int64) error
}

type SubscriptionReader interface {
	SelectByUserID(DataHandler, int64) ([]int64, error)
}

type SubscriptionDB interface {
	SubscriptionReader
	SubscriptionWriter
}

type SubscriptionSqlDB struct {
	SubscriptionReader
	SubscriptionWriter
}

func newSubscriptionDB() SubscriptionDB {
	return &SubscriptionSqlDB{
		SubscriptionWriter: &SubscriptionSqlWriter{},
		SubscriptionReader: &SubscriptionSqlReader{},
	}
}

type SubscriptionSqlReader struct{}

func (r *SubscriptionSqlReader) SelectByUserID(dh DataHandler, userID int64) ([]int64, error) {
	q := "SELECT channel_id FROM subscriptions WHERE user_id=$1"
	var result []int64
	err := sqlx.Select(dh, &result, q, userID)
	return result, dbErr(err, q)
}

type SubscriptionSqlWriter struct{}

func (w *SubscriptionSqlWriter) Create(dh DataHandler, channelID, userID int64) error {
	q := "INSERT INTO subscriptions(channel_id, user_id) VALUES($1, $2)"
	_, err := dh.Exec(q, channelID, userID)
	return dbErr(err, q)
}

func (w *SubscriptionSqlWriter) Delete(dh DataHandler, channelID, userID int64) error {
	q := "DELETE FROM subscriptions WHERE channel_id=$1 AND user_id=$2"
	_, err := dh.Exec(q, channelID, userID)
	return dbErr(err, q)
}

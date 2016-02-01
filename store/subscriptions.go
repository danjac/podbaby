package store

import (
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type SubscriptionWriter interface {
	Create(DataHandler, int64, int64) error
	Delete(DataHandler, int64, int64) error
}

type SubscriptionReader interface {
	SelectByUserID(DataHandler, *[]int64, int64) error
}

type SubscriptionStore interface {
	SubscriptionReader
	SubscriptionWriter
}

type subscriptionSqlStore struct {
	SubscriptionReader
	SubscriptionWriter
}

func newSubscriptionStore() SubscriptionStore {
	return &subscriptionSqlStore{
		SubscriptionWriter: &subscriptionSqlWriter{},
		SubscriptionReader: &subscriptionSqlReader{},
	}
}

type subscriptionSqlReader struct{}

func (r *subscriptionSqlReader) SelectByUserID(dh DataHandler, result *[]int64, userID int64) error {
	q := "SELECT channel_id FROM subscriptions WHERE user_id=$1"
	return sqlx.Select(dh, result, q, userID)
}

type subscriptionSqlWriter struct{}

func (w *subscriptionSqlWriter) Create(dh DataHandler, channelID, userID int64) error {
	q := "INSERT INTO subscriptions(channel_id, user_id) VALUES($1, $2)"
	_, err := dh.Exec(q, channelID, userID)
	return err
}

func (w *subscriptionSqlWriter) Delete(dh DataHandler, channelID, userID int64) error {
	q := "DELETE FROM subscriptions WHERE channel_id=$1 AND user_id=$2"
	_, err := dh.Exec(q, channelID, userID)
	return err
}

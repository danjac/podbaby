package store

import (
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

// SubscriptionWriter handles writes to subscription data store
type SubscriptionWriter interface {
	Create(DataHandler, int, int) error
	Delete(DataHandler, int, int) error
}

// SubscriptionReader handles reads from subscription data store
type SubscriptionReader interface {
	SelectByUserID(DataHandler, *[]int, int) error
}

// SubscriptionStore manages interactions with subscription data store
type SubscriptionStore interface {
	SubscriptionReader
	SubscriptionWriter
}

type subscriptionSQLStore struct {
	SubscriptionReader
	SubscriptionWriter
}

func newSubscriptionStore() SubscriptionStore {
	return &subscriptionSQLStore{
		SubscriptionWriter: &subscriptionSQLWriter{},
		SubscriptionReader: &subscriptionSQLReader{},
	}
}

type subscriptionSQLReader struct{}

func (r *subscriptionSQLReader) SelectByUserID(dh DataHandler, result *[]int, userID int) error {
	q := "SELECT channel_id FROM subscriptions WHERE user_id=$1"
	return handleError(sqlx.Select(dh, result, q, userID), q)
}

type subscriptionSQLWriter struct{}

func (w *subscriptionSQLWriter) Create(dh DataHandler, channelID, userID int) error {
	q := "INSERT INTO subscriptions(channel_id, user_id) VALUES($1, $2)"
	_, err := dh.Exec(q, channelID, userID)
	return handleError(err, q)
}

func (w *subscriptionSQLWriter) Delete(dh DataHandler, channelID, userID int) error {
	q := "DELETE FROM subscriptions WHERE channel_id=$1 AND user_id=$2"
	_, err := dh.Exec(q, channelID, userID)
	return handleError(err, q)
}

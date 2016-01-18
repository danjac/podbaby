package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
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

func newSubscriptionDB(db sqlx.Ext, ps purse.Purse) *SubscriptionDB {
	return &SubscriptionDB{
		SubscriptionWriter: &SubscriptionDBWriter{db, ps},
		SubscriptionReader: &SubscriptionDBReader{db, ps},
	}
}

type SubscriptionDBReader struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *SubscriptionDBReader) SelectByUserID(userID int64) ([]int64, error) {
	q, _ := db.ps.Get("select_subscriptions_by_user_id.sql")
	var result []int64
	err := sqlx.Select(db, &result, q, userID)
	return result, sqlErr(err, q)
}

type SubscriptionDBWriter struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *SubscriptionDBWriter) Create(channelID, userID int64) error {
	q, _ := db.ps.Get("insert_subscription.sql")
	_, err := db.Exec(q, channelID, userID)
	return sqlErr(err, q)
}

func (db *SubscriptionDBWriter) Delete(channelID, userID int64) error {
	q, _ := db.ps.Get("delete_subscription.sql")
	_, err := db.Exec(q, channelID, userID)
	return sqlErr(err, q)
}

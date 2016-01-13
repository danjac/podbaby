package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
)

type SubscriptionDB interface {
	Create(int64, int64) error
	Delete(int64, int64) error
	SelectByUserID(int64) ([]int64, error)
}

type defaultSubscriptionDBImpl struct {
	*sqlx.DB
	ps purse.Purse
}

func (db *defaultSubscriptionDBImpl) SelectByUserID(userID int64) ([]int64, error) {
	q, _ := db.ps.Get("select_subscriptions_by_user_id.sql")
	var result []int64
	err := db.Select(&result, q, userID)
	return result, sqlErr(err, q)
}

func (db *defaultSubscriptionDBImpl) Create(channelID, userID int64) error {
	q, _ := db.ps.Get("insert_subscription.sql")
	_, err := db.Exec(q, channelID, userID)
	return sqlErr(err, q)
}

func (db *defaultSubscriptionDBImpl) Delete(channelID, userID int64) error {
	q, _ := db.ps.Get("delete_subscription.sql")
	_, err := db.Exec(q, channelID, userID)
	return sqlErr(err, q)
}

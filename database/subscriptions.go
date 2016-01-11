package database

import (
	"github.com/danjac/podbaby/sql"
	"github.com/jmoiron/sqlx"
)

type SubscriptionDB interface {
	Create(int64, int64) error
	Delete(int64, int64) error
	SelectByUserID(int64) ([]int64, error)
}

type defaultSubscriptionDBImpl struct {
	*sqlx.DB
}

func (db *defaultSubscriptionDBImpl) SelectByUserID(userID int64) ([]int64, error) {
	q, _ := sql.Queries.Get("select_subscriptions_by_user_id.sql")
	var result []int64
	err := db.Select(&result, q, userID)
	return result, err
}

func (db *defaultSubscriptionDBImpl) Create(channelID, userID int64) error {
	q, _ := sql.Queries.Get("insert_subscription.sql")
	_, err := db.Exec(q, channelID, userID)
	return err
}

func (db *defaultSubscriptionDBImpl) Delete(channelID, userID int64) error {
	q, _ := sql.Queries.Get("delete_subscription.sql")
	_, err := db.Exec(q, channelID, userID)
	return err
}

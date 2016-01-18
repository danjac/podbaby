package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
)

type UserWriter interface {
	Create(*models.User) error
	UpdateEmail(string, int64) error
	UpdatePassword(string, int64) error
	DeleteUser(int64) error
}

type UserReader interface {
	GetByID(int64) (*models.User, error)
	GetByNameOrEmail(string) (*models.User, error)
	IsName(string) (bool, error)
	IsEmail(string, int64) (bool, error)
}

func newUserDB(db sqlx.Ext, ps purse.Purse) *UserDB {
	return &UserDB{
		UserWriter: &UserDBWriter{db, ps},
		UserReader: &UserDBReader{db, ps},
	}
}

// UserDB handles all user queries
type UserDB struct {
	sqlx.Ext
	UserReader
	UserWriter
}

type UserDBWriter struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *UserDBWriter) UpdateEmail(email string, userID int64) error {
	q, _ := db.ps.Get("update_user_email.sql")
	_, err := db.Exec(q, email, userID)
	return sqlErr(err, q)
}

func (db *UserDBWriter) UpdatePassword(password string, userID int64) error {
	q, _ := db.ps.Get("update_user_password.sql")
	_, err := db.Exec(q, password, userID)
	return sqlErr(err, q)
}
func (db *UserDBWriter) Create(user *models.User) error {
	q, _ := db.ps.Get("insert_user.sql")
	q, args, err := sqlx.Named(q, user)
	if err != nil {
		return sqlErr(err, q)
	}
	return sqlErr(db.QueryRowx(db.Rebind(q), args...).Scan(&user.ID), q)
}

func (db *UserDBWriter) DeleteUser(userID int64) error {
	q, _ := db.ps.Get("delete_user.sql")
	_, err := db.Exec(q, userID)
	return sqlErr(err, q)
}

type UserDBReader struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *UserDBReader) GetByID(id int64) (*models.User, error) {
	q, _ := db.ps.Get("get_user_by_id.sql")
	user := &models.User{}
	err := sqlx.Get(db, user, q, id)
	return user, sqlErr(err, q)
}

func (db *UserDBReader) GetByNameOrEmail(identifier string) (*models.User, error) {
	q, _ := db.ps.Get("get_user_by_name_or_email.sql")
	user := &models.User{}
	err := sqlx.Get(db, user, q, identifier)
	return user, sqlErr(err, q)
}

func (db *UserDBReader) IsName(name string) (bool, error) {
	q, _ := db.ps.Get("user_name_exists.sql")
	var count int64
	if err := db.QueryRowx(q, name).Scan(&count); err != nil {
		return false, sqlErr(err, q)
	}
	return count > 0, nil

}

func (db *UserDBReader) IsEmail(email string, userID int64) (bool, error) {

	qname := "user_email_exists.sql"
	args := []interface{}{email}

	if userID != 0 {
		qname = "user_email_exists_with_id.sql"
		args = append(args, userID)
	}

	q, _ := db.ps.Get(qname)

	var count int64

	if err := db.QueryRowx(q, args...).Scan(&count); err != nil {
		return false, sqlErr(err, q)
	}
	return count > 0, nil
}

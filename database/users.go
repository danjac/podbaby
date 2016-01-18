package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
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

func newUserDB(db sqlx.Ext) *UserDB {
	return &UserDB{
		UserWriter: &UserDBWriter{db},
		UserReader: &UserDBReader{db},
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
}

func (db *UserDBWriter) UpdateEmail(email string, userID int64) error {
	q := "UPDATE users SET email=$1 WHERE id=$2"
	_, err := db.Exec(q, email, userID)
	return dbErr(err, q)
}

func (db *UserDBWriter) UpdatePassword(password string, userID int64) error {
	q := "UPDATE users SET password=$1 WHERE id=$2"
	_, err := db.Exec(q, password, userID)
	return dbErr(err, q)
}
func (db *UserDBWriter) Create(user *models.User) error {
	q := `INSERT INTO users(name, email, password)
VALUES (:name, :email, :password) RETURNING id`
	q, args, err := sqlx.Named(q, user)
	if err != nil {
		return dbErr(err, q)
	}
	return dbErr(db.QueryRowx(db.Rebind(q), args...).Scan(&user.ID), q)
}

func (db *UserDBWriter) DeleteUser(userID int64) error {
	q := "DELETE FROM users WHERE id=$1"
	_, err := db.Exec(q, userID)
	return dbErr(err, q)
}

type UserDBReader struct {
	sqlx.Ext
}

func (db *UserDBReader) GetByID(id int64) (*models.User, error) {
	q := "SELECT * FROM users WHERE id=$1"
	user := &models.User{}
	err := sqlx.Get(db, user, q, id)
	return user, dbErr(err, q)
}

func (db *UserDBReader) GetByNameOrEmail(identifier string) (*models.User, error) {
	q := "SELECT * FROM users WHERE email=$1 or name=$1"
	user := &models.User{}
	err := sqlx.Get(db, user, q, identifier)
	return user, dbErr(err, q)
}

func (db *UserDBReader) IsName(name string) (bool, error) {
	q := "SELECT COUNT(id) FROM users WHERE name=$1"
	var count int64
	if err := db.QueryRowx(q, name).Scan(&count); err != nil {
		return false, dbErr(err, q)
	}
	return count > 0, nil

}

func (db *UserDBReader) IsEmail(email string, userID int64) (bool, error) {

	q := "SELECT COUNT(id) FROM users WHERE email ILIKE $1"
	args := []interface{}{email}

	if userID != 0 {
		q += " AND id != $2"
		args = append(args, userID)
	}

	var count int64

	if err := db.QueryRowx(q, args...).Scan(&count); err != nil {
		return false, dbErr(err, q)
	}
	return count > 0, nil
}

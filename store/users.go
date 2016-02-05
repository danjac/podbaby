package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

// UserWriter does writes to user data store
type UserWriter interface {
	// Create adds a new user to the data store
	Create(DataHandler, *models.User) error
	// UpdateEmail updates email address in the data store
	UpdateEmail(DataHandler, string, int) error
	// UpdatePassword updates password in the data store
	UpdatePassword(DataHandler, string, int) error
	// DeleteUser removes user from data store
	DeleteUser(DataHandler, int) error
}

// UserReader does reads from user data store
type UserReader interface {
	// GetbYID retrieves user by ID from the data store
	GetByID(DataHandler, *models.User, int) error
	// GetNameOrEmail retrieves user by name or email from the data store
	GetByNameOrEmail(DataHandler, *models.User, string) error
	// IsName checks if a user name exists in the data store
	IsName(DataHandler, string) (bool, error)
	// IsEmail checks if an email exists in the data store
	IsEmail(DataHandler, string, int) (bool, error)
}

// UserStore manages the user data store
type UserStore interface {
	UserReader
	UserWriter
}

// UserStore handles all user queries
type userSQLStore struct {
	UserReader
	UserWriter
}

func newUserStore() UserStore {
	return &userSQLStore{
		UserWriter: &userSQLWriter{},
		UserReader: &userSQLReader{},
	}
}

type userSQLWriter struct{}

func (w *userSQLWriter) UpdateEmail(dh DataHandler, email string, userID int) error {
	q := "UPDATE users SET email=$1 WHERE id=$2"
	_, err := dh.Exec(q, email, userID)
	return handleError(err, q)
}

func (w *userSQLWriter) UpdatePassword(dh DataHandler, password string, userID int) error {
	q := "UPDATE users SET password=$1 WHERE id=$2"
	_, err := dh.Exec(q, password, userID)
	return handleError(err, q)
}
func (w *userSQLWriter) Create(dh DataHandler, user *models.User) error {
	q := "INSERT INTO users(name, email, password) VALUES (:name, :email, :password) RETURNING id"
	q, args, err := sqlx.Named(q, user)
	if err != nil {
		return handleError(err, q)
	}
	return handleError(dh.QueryRowx(dh.Rebind(q), args...).Scan(&user.ID), q)
}

func (w *userSQLWriter) DeleteUser(dh DataHandler, userID int) error {
	q := "DELETE FROM users WHERE id=$1"
	_, err := dh.Exec(q, userID)
	return handleError(err, q)
}

type userSQLReader struct{}

func (r *userSQLReader) GetByID(dh DataHandler, user *models.User, id int) error {
	q := "SELECT * FROM users WHERE id=$1"
	return handleError(sqlx.Get(dh, user, q, id), q)
}

func (r *userSQLReader) GetByNameOrEmail(dh DataHandler, user *models.User, identifier string) error {
	q := "SELECT * FROM users WHERE email=$1 or name=$1"
	return handleError(sqlx.Get(dh, user, q, identifier), q)
}

func (r *userSQLReader) IsName(dh DataHandler, name string) (bool, error) {
	q := "SELECT COUNT(id) FROM users WHERE name=$1"
	var count int
	if err := dh.QueryRowx(q, name).Scan(&count); err != nil {
		return false, handleError(err, q)
	}
	return count > 0, nil

}

func (r *userSQLReader) IsEmail(dh DataHandler, email string, userID int) (bool, error) {

	q := "SELECT COUNT(id) FROM users WHERE email ILIKE $1"
	args := []interface{}{email}

	if userID != 0 {
		q += " AND id != $2"
		args = append(args, userID)
	}

	var count int

	if err := dh.QueryRowx(q, args...).Scan(&count); err != nil {
		return false, handleError(err, q)
	}
	return count > 0, nil
}

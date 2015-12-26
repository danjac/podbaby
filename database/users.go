package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type UserDB interface {
	GetByID(int64) (*models.User, error)
	GetByNameOrEmail(string) (*models.User, error)
	Create(*models.User) error
	IsName(string) (bool, error)
	IsEmail(string) (bool, error)
}

type defaultUserDBImpl struct {
	*sqlx.DB
}

func (db *defaultUserDBImpl) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := db.Get(user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func (db *defaultUserDBImpl) GetByNameOrEmail(identifier string) (*models.User, error) {
	user := &models.User{}
	err := db.Get(user, "SELECT * FROM users WHERE email=$1 or name=$1", identifier)
	return user, err
}

func (db *defaultUserDBImpl) Create(user *models.User) error {
	query, args, err := sqlx.Named(`INSERT INTO users(name, email, password) 
    VALUES (:name, :email, :password) RETURNING id`, user)
	if err != nil {
		return err
	}
	return db.QueryRow(db.Rebind(query), args...).Scan(&user.ID)
}

func (db *defaultUserDBImpl) IsName(name string) (bool, error) {
	var count int64
	if err := db.QueryRow("SELECT COUNT(id) FROM users WHERE name=$1", name).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil

}

func (db *defaultUserDBImpl) IsEmail(email string) (bool, error) {
	var count int64
	if err := db.QueryRow("SELECT COUNT(id) FROM users WHERE email=$1", email).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

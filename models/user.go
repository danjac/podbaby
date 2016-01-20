package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	Bookmarks     []int64   `db:"-" json:"bookmarks"`
	Subscriptions []int64   `db:"-" json:"subscriptions"`
	Plays         []Play    `db:"-" json:"plays"`
}

func (user *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	if user.Password == "" {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

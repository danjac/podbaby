package models

import (
	"github.com/danjac/podbaby/models/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
	"time"
)

// User is the current authenticated user
type User struct {
	ID            int       `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	Bookmarks     []int     `db:"-" json:"bookmarks"`
	Subscriptions []int     `db:"-" json:"subscriptions"`
	Plays         []Play    `db:"-" json:"plays"`
}

// SetPassword sets the Password field with an encrypted password
func (user *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

// CheckPassword checks if the password matches the encrypted password
func (user *User) CheckPassword(password string) bool {
	if user.Password == "" {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

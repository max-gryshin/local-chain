package models

import (
	"golang.org/x/crypto/bcrypt"

	"time"
)

const (
	StateActive   = 1
	StateInActive = 2
	StateBlocked  = 3
	StateDeleted  = 4
)

type Users []*User

type User struct {
	ID         int       `json:"id"          db:"id"`
	Email      string    `json:"email"       db:"email"`
	Password   string    `json:"password"    db:"password_hash"`
	FirstName  *string   `json:"first_name"  db:"first_name"`
	LastName   *string   `json:"last_name"   db:"last_name"`
	MiddleName *string   `json:"middle_name" db:"middle_name"`
	Status     int       `json:"status"      db:"status"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"     goqu:"skipupdate"`
	UpdatedAt  time.Time `json:"updated_at"  db:"updated_at"`
	CreatedBy  int       `json:"created_by"  db:"created_by"`
	UpdatedBy  int       `json:"updated_by"  db:"updated_by"`
	Roles      string    `json:"roles"       db:"roles"`
}

// SetPassword sets a new password stored as hash.
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)

	return nil
}

// InvalidPassword returns true if the given password does not match the hash.
func (u *User) InvalidPassword(password string) bool {
	if u.Password == "" && password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err != nil
}

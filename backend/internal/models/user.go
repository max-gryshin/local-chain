package models

import (
	"golang.org/x/crypto/bcrypt"

	"time"
)

//todo: move to profile
const (
	StateHalfRegistration = 1
	StateRegistration     = 2
	StateActive           = 3
	StateBlocked          = 4
	StateDeleted          = 5
)

type Users []*User

var userFields = map[string][]string{
	"get":    {"id", "state", "created_at", "email"},
	"update": {"user_name", "state", "email"},
}

type User struct {
	ID         int       `json:"id"          db:"id"`
	Email      string    `json:"email"       db:"email"`
	Password   string    `json:"password"    db:"password_hash"`
	FirstName  *string   `json:"first_name"  db:"first_name"`
	LastName   *string   `json:"last_name"   db:"last_name"`
	MiddleName *string   `json:"middle_name" db:"middle_name"`
	Status     int       `json:"status"      db:"status"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"  db:"updated_at"`
	CreatedBy  int       `json:"created_by"  db:"created_by"`
	UpdatedBy  int       `json:"updated_by"  db:"updated_by"`
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

func GetAllowedUserFieldsByMethod(method string) []string {
	return userFields[method]
}

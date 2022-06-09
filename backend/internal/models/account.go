package models

import "time"

type Accounts []*Account

type Account struct {
	ID        int       `json:"id"          db:"id"`
	Phone     string    `json:"phone"       db:"phone"`
	Dob       time.Time `json:"dob"         db:"dob"`
	Status    int       `json:"status"      db:"status"`
	UserID    int       `json:"user_id"     db:"user_id"`
	CreatedAt time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt time.Time `json:"updated_at"  db:"updated_at"`
	CreatedBy int       `json:"created_by"  db:"created_by"`
	UpdatedBy int       `json:"updated_by"  db:"updated_by"`
}

package models

import "time"

type Orders []*Order

type Order struct {
	ID          int     `json:"id"              db:"id"`
	Status      int     `json:"status"          db:"status"`
	Amount      float64 `json:"amount"          db:"amount"`
	Description string  `json:"description"     db:"description"`
	// todo: change postgres type to enum
	RequestReasons []string  `json:"request_reasons" db:"request_reasons"`
	WalletID       int       `json:"wallet_id"       db:"wallet_id"`
	CreatedAt      time.Time `json:"created_at"      db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"      db:"updated_at"`
	CreatedBy      int       `json:"created_by"      db:"created_by"`
	UpdatedBy      int       `json:"updated_by"      db:"updated_by"`
}

package models

import "time"

type Wallets []*Wallet

type Wallet struct {
	ID         int       `json:"id"          db:"id"`
	Status     int       `json:"status"      db:"status"`
	WalletID   string    `json:"description" db:"description"`
	PrivateKey string    `json:"wallet_id"   db:"wallet_id"`
	AccountID  int       `json:"account_id"  db:"account_id"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"  db:"updated_at"`
	CreatedBy  int       `json:"created_by"  db:"created_by"`
	UpdatedBy  int       `json:"updated_by"  db:"updated_by"`
}

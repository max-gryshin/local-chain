package dto

import (
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type Wallets []*Wallet

type Wallet struct {
	ID         int       `json:"id" validate:"required"`
	Status     int       `json:"status"`
	WalletID   string    `json:"wallet_id"`
	PrivateKey string    `json:"private_key"`
	AccountID  int       `json:"account_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  int       `json:"created_by"`
	UpdatedBy  int       `json:"updated_by"`
}

func LoadWalletDTOFromModel(model *models.Wallet) *Wallet {
	return &Wallet{
		ID:        model.ID,
		Status:    model.Status,
		WalletID:  model.WalletID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
	}
}

func LoadWalletModelFromDTO(dto *Wallet) *models.Wallet {
	return &models.Wallet{
		ID:         dto.ID,
		Status:     dto.Status,
		WalletID:   dto.WalletID,
		PrivateKey: dto.PrivateKey,
		AccountID:  dto.AccountID,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
		CreatedBy:  dto.CreatedBy,
		UpdatedBy:  dto.UpdatedBy,
	}
}

func LoadWalletDTOCollectionFromModel(walletsModel models.Wallets) Wallets {
	var walletsDTO Wallets
	for _, wallet := range walletsModel {
		walletsDTO = append(walletsDTO, LoadWalletDTOFromModel(wallet))
	}
	return walletsDTO
}

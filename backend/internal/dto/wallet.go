package dto

import (
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type Wallets []*Wallet

type WalletCreate struct {
	Status     int    `json:"status"`
	WalletID   string `json:"wallet_id"`
	PrivateKey string `json:"private_key"`
	AccountID  int    `json:"account_id"`
}

type Wallet struct {
	ID        int       `json:"id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	WalletCreate
}

func LoadWalletDTOFromModel(model *models.Wallet) *Wallet {
	return &Wallet{
		ID: model.ID,
		WalletCreate: WalletCreate{
			Status:     model.Status,
			WalletID:   model.WalletID,
			PrivateKey: model.PrivateKey,
			AccountID:  model.AccountID,
		},
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
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func LoadWalletDTOCollectionFromModel(walletsModel models.Wallets) Wallets {
	var walletsDTO Wallets
	for _, wallet := range walletsModel {
		walletsDTO = append(walletsDTO, LoadWalletDTOFromModel(wallet))
	}
	return walletsDTO
}

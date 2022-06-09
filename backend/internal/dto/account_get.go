package dto

import (
	"time"

	"github.com/max-gryshin/local-chain/internal/models"
)

type GetAccounts []*GetAccount

type GetAccount struct {
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int       `json:"user_id"`
	UpdatedBy int       `json:"updated_by"`
	AccountOwnerUpdate
}

func LoadGetAccountDTOFromModel(model *models.Account) *GetAccount {
	return &GetAccount{
		Status:             model.Status,
		UserID:             model.UserID,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		CreatedBy:          model.CreatedBy,
		UpdatedBy:          model.UpdatedBy,
		AccountOwnerUpdate: *LoadAccountOwnerUpdateDTOFromModel(model),
	}
}

func LoadAccountModelFromGetAccountDTO(dto *GetAccount) *models.Account {
	return &models.Account{
		ID:        dto.ID,
		Phone:     dto.Phone,
		Dob:       dto.Dob,
		Status:    dto.Status,
		UserID:    dto.UserID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}
}

func LoadGetAccountsDTOCollectionFromModel(accountsModel models.Accounts) GetAccounts {
	var accountsDTO GetAccounts
	for _, account := range accountsModel {
		accountsDTO = append(accountsDTO, LoadGetAccountDTOFromModel(account))
	}
	return accountsDTO
}

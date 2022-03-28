package dto

import (
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type Accounts []*Account

type Account struct {
	ID        int       `json:"id" validate:"required"`
	Phone     string    `json:"phone"`
	Dob       time.Time `json:"dob"`
	Status    int       `json:"status"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
}

func LoadAccountDTOFromModel(model *models.Account) *Account {
	return &Account{
		ID:        model.ID,
		Phone:     model.Phone,
		Dob:       model.Dob,
		Status:    model.Status,
		UserID:    model.UserID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
	}
}

func LoadAccountModelFromDTO(dto *Account) *models.Account {
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

func LoadAccountDTOCollectionFromModel(accountsModel models.Accounts) Accounts {
	var accountsDTO Accounts
	for _, account := range accountsModel {
		accountsDTO = append(accountsDTO, LoadAccountDTOFromModel(account))
	}
	return accountsDTO
}

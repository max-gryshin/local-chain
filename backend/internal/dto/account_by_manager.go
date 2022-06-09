package dto

import (
	"time"

	"github.com/max-gryshin/local-chain/internal/models"
)

type AccountByManager struct {
	Status int `json:"status"`
	UserID int `json:"user_id"`
	AccountOwnerUpdateRequest
}

func LoadAccountByManagerModelFromDTO(dto *AccountByManager) *models.Account {
	return &models.Account{
		Phone:     dto.Phone,
		Dob:       dto.Dob,
		Status:    dto.Status,
		UpdatedAt: time.Now(),
	}
}

func LoadAccountByManagerDTOFromModel(model *models.Account) *AccountByManager {
	return &AccountByManager{
		Status: model.Status,
		UserID: model.UserID,
		AccountOwnerUpdateRequest: AccountOwnerUpdateRequest{
			Dob:   model.Dob,
			Phone: model.Phone,
		},
	}
}

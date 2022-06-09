package dto

import (
	"time"

	"github.com/max-gryshin/local-chain/internal/models"
)

type AccountOwnerUpdateRequest struct {
	Dob   time.Time `json:"dob"`
	Phone string    `json:"phone"`
}

type AccountOwnerUpdate struct {
	ID int `json:"id"         validate:"required"`
	AccountOwnerUpdateRequest
}

func LoadAccountOwnerUpdateModelFromDTO(dto *AccountOwnerUpdate) *models.Account {
	return &models.Account{
		ID:        dto.ID,
		Phone:     dto.Phone,
		Dob:       dto.Dob,
		UpdatedAt: time.Now(),
	}
}

func LoadAccountOwnerUpdateDTOFromModel(model *models.Account) *AccountOwnerUpdate {
	return &AccountOwnerUpdate{
		ID: model.ID,
		AccountOwnerUpdateRequest: AccountOwnerUpdateRequest{
			Dob:   model.Dob,
			Phone: model.Phone,
		},
	}
}

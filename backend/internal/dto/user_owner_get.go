package dto

import (
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type GetUserOwners []*GetUserOwner

type GetUserOwner struct {
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	UpdateUserOwner
}

func LoadGetUserOwnerDTOFromModel(model *models.User) *GetUserOwner {
	return &GetUserOwner{
		Status:          model.Status,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		CreatedBy:       model.CreatedBy,
		UpdatedBy:       model.UpdatedBy,
		UpdateUserOwner: *LoadUpdateUserOwnerDTOFromModel(model),
	}
}

func LoadUserModelFromGetUserOwnerDTO(dto *GetUserOwner) *models.User {
	return &models.User{
		ID:         dto.ID,
		Email:      *dto.Email,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		MiddleName: dto.MiddleName,
		Status:     dto.Status,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
		CreatedBy:  dto.CreatedBy,
		UpdatedBy:  dto.UpdatedBy,
	}
}

func LoadGetUserOwnerDTOCollectionFromModel(usersModel models.Users) GetUserOwners {
	var usersDTO GetUserOwners
	for _, user := range usersModel {
		usersDTO = append(usersDTO, LoadGetUserOwnerDTOFromModel(user))
	}
	return usersDTO
}

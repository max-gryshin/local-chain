package dto

import (
	"time"

	"github.com/max-gryshin/local-chain/internal/models"
)

type UsersByManager []*UserByManager

type UserByManager struct {
	Roles []string `json:"roles"`
	GetUserOwner
}

func LoadUserByManagerDTOFromModel(model *models.User) *UserByManager {
	return &UserByManager{
		Roles:        model.Roles,
		GetUserOwner: *LoadGetUserOwnerDTOFromModel(model),
	}
}

func LoadUserModelFromUserByManagerDTO(dto *UserByManager) *models.User {
	return &models.User{
		ID:         dto.ID,
		Email:      *dto.Email,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		MiddleName: dto.MiddleName,
		Status:     dto.Status,
		Roles:      dto.Roles,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  time.Now(),
		CreatedBy:  dto.CreatedBy,
		UpdatedBy:  dto.UpdatedBy,
		ManagerID:  dto.ManagerID,
	}
}

func LoadUsersByManagerDTOCollectionFromModel(usersModel models.Users) UsersByManager {
	var usersDTO UsersByManager
	for _, user := range usersModel {
		usersDTO = append(usersDTO, LoadUserByManagerDTOFromModel(user))
	}
	return usersDTO
}

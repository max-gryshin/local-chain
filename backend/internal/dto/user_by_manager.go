package dto

import (
	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type UsersByManager []*UserByManager

type UserByManager struct {
	Password string `json:"password" validate:"gte=6,lte=50"`
	Roles    string `json:"roles"    validate:"json"`
	GetUserOwner
}

func LoadUserByManagerDTOFromModel(model *models.User) *UserByManager {
	return &UserByManager{
		Password:     model.Password,
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
		Password:   dto.Password,
		Roles:      dto.Roles,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
		CreatedBy:  dto.CreatedBy,
		UpdatedBy:  dto.UpdatedBy,
	}
}

func LoadUsersByManagerDTOCollectionFromModel(usersModel models.Users) UsersByManager {
	var usersDTO UsersByManager
	for _, user := range usersModel {
		usersDTO = append(usersDTO, LoadUserByManagerDTOFromModel(user))
	}
	return usersDTO
}

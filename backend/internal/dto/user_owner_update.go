package dto

import (
	"github.com/max-gryshin/local-chain/internal/models"
)

type UpdateUserOwners []*UpdateUserOwner

type UpdateUserOwnerRequest struct {
	Email      *string `json:"email"       validate:"required,email"`
	FirstName  *string `json:"first_name"  validate:"alpha"`
	LastName   *string `json:"last_name"   validate:"alpha"`
	MiddleName *string `json:"middle_name" validate:"alpha"`
}

type UpdateUserOwner struct {
	ID int `json:"id"          validate:"required"`
	UpdateUserOwnerRequest
}

func LoadUpdateUserOwnerDTOFromModel(model *models.User) *UpdateUserOwner {
	return &UpdateUserOwner{
		ID: model.ID,
		UpdateUserOwnerRequest: UpdateUserOwnerRequest{
			Email:      &model.Email,
			FirstName:  model.FirstName,
			LastName:   model.LastName,
			MiddleName: model.MiddleName,
		},
	}
}

func LoadUserModelFromUpdateUserOwnerDTO(dto *UpdateUserOwner) *models.User {
	return &models.User{
		ID:         dto.ID,
		Email:      *dto.Email,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		MiddleName: dto.MiddleName,
	}
}

func LoadUpdateUserDTOCollectionFromModel(usersModel models.Users) UpdateUserOwners {
	var usersDTO UpdateUserOwners
	for _, user := range usersModel {
		usersDTO = append(usersDTO, LoadUpdateUserOwnerDTOFromModel(user))
	}
	return usersDTO
}

package dto

import (
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type Users []*User

type User struct {
	ID        int       `json:"id" validate:"required"`
	Password  string    `json:"password"   ` //validate:"gte=6,lte=50"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"      validate:"email"`
}

func LoadUserDTOFromModel(model *models.User) *User {
	return &User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
	}
}

func LoadUserModelFromDTO(dto *User) *models.User {
	return &models.User{
		ID:        dto.ID,
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: dto.CreatedAt,
	}
}

func LoadUserDTOCollectionFromModel(usersModel models.Users) Users {
	var usersDTO Users
	for _, user := range usersModel {
		usersDTO = append(usersDTO, LoadUserDTOFromModel(user))
	}
	return usersDTO
}

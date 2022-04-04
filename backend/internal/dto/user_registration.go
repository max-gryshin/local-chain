package dto

import (
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type UserRegistration struct {
	Status   int      `json:"status"   validate:"required,gte=1,lte=4"`
	Password string   `json:"password" validate:"required,gte=6,lte=50"`
	Roles    []string `json:"roles"`
	UpdateUserOwnerRequest
}

func LoadUserModelFromUserRegistrationDTO(dto *UserRegistration) *models.User {
	return &models.User{
		Email:      *dto.Email,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		MiddleName: dto.MiddleName,
		Status:     dto.Status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

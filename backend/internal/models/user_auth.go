package models

type Auth struct {
	Email    string `json:"email" validate:"required,gte=3,lte=50"`
	Password string `json:"password" validate:"required,gte=6,lte=50"`
}

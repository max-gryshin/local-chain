package dto

import "github.com/ZmaximillianZ/local-chain/internal/models"

type OrderByManager struct {
	Status int `json:"status"`
}

func LoadOrderByManagerDTOFromModel(model *models.Order) *OrderByManager {
	return &OrderByManager{Status: model.Status}
}

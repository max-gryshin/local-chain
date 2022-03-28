package dto

import (
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/models"
)

type Orders []*Order

type Order struct {
	ID            int       `json:"id" validate:"required"`
	Status        int       `json:"status"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	RequestReason []string  `json:"request_reason"`
	WalletID      int       `json:"wallet_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     int       `json:"created_by"`
	UpdatedBy     int       `json:"updated_by"`
}

func LoadOrderDTOFromModel(model *models.Order) *Order {
	return &Order{
		ID:            model.ID,
		Status:        model.Status,
		Amount:        model.Amount,
		Description:   model.Description,
		RequestReason: model.RequestReason,
		WalletID:      model.WalletID,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		CreatedBy:     model.CreatedBy,
		UpdatedBy:     model.UpdatedBy,
	}
}

func LoadOrderModelFromDTO(dto *Order) *models.Order {
	return &models.Order{
		ID:            dto.ID,
		Status:        dto.Status,
		Amount:        dto.Amount,
		Description:   dto.Description,
		RequestReason: dto.RequestReason,
		WalletID:      dto.WalletID,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
		CreatedBy:     dto.CreatedBy,
		UpdatedBy:     dto.UpdatedBy,
	}
}

func LoadOrderDTOCollectionFromModel(ordersModel models.Orders) Orders {
	var ordersDTO Orders
	for _, order := range ordersModel {
		ordersDTO = append(ordersDTO, LoadOrderDTOFromModel(order))
	}
	return ordersDTO
}

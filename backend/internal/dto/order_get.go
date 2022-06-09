package dto

import (
	"time"

	"github.com/max-gryshin/local-chain/internal/models"
)

type GetOrders []*GetOrder

type OrderRequest struct {
	Amount        float64  `json:"amount"`
	Description   string   `json:"description"`
	RequestReason []string `json:"request_reason"`
	WalletID      int      `json:"wallet_id"`
}

type GetOrder struct {
	ID        int       `json:"id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	OrderRequest
	OrderByManager
}

func LoadGetOrderDTOFromModel(model *models.Order) *GetOrder {
	return &GetOrder{
		ID: model.ID,
		OrderRequest: OrderRequest{
			Amount:        model.Amount,
			Description:   model.Description,
			RequestReason: model.RequestReasons,
			WalletID:      model.WalletID,
		},
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		CreatedBy:      model.CreatedBy,
		UpdatedBy:      model.UpdatedBy,
		OrderByManager: *LoadOrderByManagerDTOFromModel(model),
	}
}

func LoadGetOrderModelFromDTO(dto *GetOrder) *models.Order {
	return &models.Order{
		ID:             dto.ID,
		Status:         dto.Status,
		Amount:         dto.Amount,
		Description:    dto.Description,
		RequestReasons: dto.RequestReason,
		WalletID:       dto.WalletID,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
		CreatedBy:      dto.CreatedBy,
		UpdatedBy:      dto.UpdatedBy,
	}
}

func LoadGetOrdersDTOCollectionFromModel(ordersModel models.Orders) GetOrders {
	var ordersDTO GetOrders
	for _, order := range ordersModel {
		ordersDTO = append(ordersDTO, LoadGetOrderDTOFromModel(order))
	}
	return ordersDTO
}

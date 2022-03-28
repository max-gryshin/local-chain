package controllers

import (
	"net/http"

	"github.com/ZmaximillianZ/local-chain/internal/contractions"
	"github.com/ZmaximillianZ/local-chain/internal/dto"
	"github.com/ZmaximillianZ/local-chain/internal/e"
	"github.com/ZmaximillianZ/local-chain/internal/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// OrderController is HTTP controller for manage order
type OrderController struct {
	repo         contractions.OrderRepository
	errorHandler e.ErrorHandler
	BaseController
}

// NewOrderController return new instance of OrderController
func NewOrderController(repo contractions.OrderRepository, errorHandler e.ErrorHandler, v *validator.Validate) *OrderController {
	return &OrderController{
		repo:           repo,
		errorHandler:   errorHandler,
		BaseController: BaseController{*v},
	}
}

// GetByID example: /api/order/{id}
func (ctr *OrderController) GetByID(c echo.Context) error {
	var (
		err   error
		order models.Order
	)
	if order, err = ctr.getOrderByID(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.LoadOrderDTOFromModel(&order))
}

// GetOrders return list of users
// example: /api/order
func (ctr *OrderController) GetOrders(c echo.Context) error {
	var (
		orders models.Orders
		err    error
	)
	if orders, err = ctr.repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadOrderDTOCollectionFromModel(orders))
}

func (ctr *OrderController) getOrderByID(c echo.Context) (models.Order, error) {
	var (
		id    int64
		err   error
		order models.Order
	)
	if id, err = ctr.BaseController.GetID(c); err != nil {
		return order, err
	}
	if order, err = ctr.repo.GetByID(int(id)); err != nil {
		return order, err
	}

	return order, err
}

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

// GetByID       godoc
// @Summary      get order
// @Description  get order by id
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "GetOrder ID"
// @Success      200  {object}  dto.GetOrder
// @Security     ApiKeyAuth
// @Router       /api/order/{id} [get]
func (ctr *OrderController) GetByID(c echo.Context) error {
	var (
		err   error
		order models.Order
	)
	if order, err = ctr.getOrderByID(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.LoadGetOrderDTOFromModel(&order))
}

// GetOrders     godoc
// @Summary      get orders
// @Description  get orders
// @Tags         order
// @Accept       json
// @Produce      json
// @Success      200  {object} dto.GetOrders
// @Security     ApiKeyAuth
// @Router       /api/order [get]
func (ctr *OrderController) GetOrders(c echo.Context) error {
	var (
		orders models.Orders
		err    error
	)
	if orders, err = ctr.repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadGetOrdersDTOCollectionFromModel(orders))
}

// CashOut       godoc
// @Summary      cash out
// @Description  create an order to cash out money
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        message  body  dto.OrderRequest  true  "OrderRequest"
// @Success      200  {object}  dto.GetOrder
// @Security     ApiKeyAuth
// @Router       /api/account/{id}/cash-out [post]
// todo: create dto
func (ctr *OrderController) CashOut(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
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

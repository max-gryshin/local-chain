package controllers

import (
	"net/http"

	"github.com/ZmaximillianZ/local-chain/internal/e"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// TransactionController is HTTP controller for manage transaction
type TransactionController struct {
	errorHandler e.ErrorHandler
	BaseController
}

// NewTransactionController return new instance of TransactionController
func NewTransactionController(errorHandler e.ErrorHandler, v *validator.Validate) *TransactionController {
	return &TransactionController{
		errorHandler:   errorHandler,
		BaseController: BaseController{*v},
	}
}

// GetTransactions godoc
// @Summary        get transactions
// @Description    getting all transaction history
// @Tags           transaction
// @Accept         json
// @Produce        json
// @Success        200  {object} dto.GetOrders
// @Security       ApiKeyAuth
// @Router         /api/transaction [get]
// todo: create dto
func (ctr *TransactionController) GetTransactions(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "2", "3"})
}

// GetUserTransactions godoc
// @Summary            get user transactions
// @Description        getting transactions certain user
// @Tags               transaction
// @Accept             json
// @Produce            json
// @Success            200  {object} dto.GetOrders
// @Security           ApiKeyAuth
// @Router             /api/transaction/user/{id} [get]
// todo: create dto
func (ctr *TransactionController) GetUserTransactions(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "2", "3"})
}

// SendTransaction  godoc
// @Summary         send transaction
// @Description     sending transaction
// @Tags            transaction
// @Accept          json
// @Produce         json
// @Param           message  body  dto.Wallet  true  "Wallet"
// @Success         200  {object} dto.Wallet
// @Security        ApiKeyAuth
// @Router          /api/transaction [post]
// todo: create dto
func (ctr *TransactionController) SendTransaction(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

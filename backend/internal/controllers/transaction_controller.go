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

// GetTransactions example: /api/transaction/all
func (ctr *TransactionController) GetTransactions(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "2", "3"})
}

// GetUserTransactions example: /api/transaction/{userId}/all
func (ctr *TransactionController) GetUserTransactions(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "2", "3"})
}

// SendTransaction example: /api/transaction
func (ctr *TransactionController) SendTransaction(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

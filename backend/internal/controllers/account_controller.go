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

// AccountController is HTTP controller for manage account
type AccountController struct {
	repo         contractions.AccountRepository
	errorHandler e.ErrorHandler
	BaseController
}

// NewAccountController return new instance of AccountController
func NewAccountController(repo contractions.AccountRepository, errorHandler e.ErrorHandler, v *validator.Validate) *AccountController {
	return &AccountController{
		repo:           repo,
		errorHandler:   errorHandler,
		BaseController: BaseController{*v},
	}
}

// GetByID example: /api/account/{id}/
func (ctr *AccountController) GetByID(c echo.Context) error {
	var (
		err     error
		account models.Account
	)
	if account, err = ctr.getAccountByID(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.LoadAccountDTOFromModel(&account))
}

// GetAccounts return list of accounts
// example: /api/account/all
func (ctr *AccountController) GetAccounts(c echo.Context) error {
	var (
		accounts models.Accounts
		err      error
	)
	if accounts, err = ctr.repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadAccountDTOCollectionFromModel(accounts))
}

// CashOut example: /api/account/{id}/cash-out
func (ctr *AccountController) CashOut(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func (ctr *AccountController) getAccountByID(c echo.Context) (models.Account, error) {
	var (
		id      int64
		err     error
		account models.Account
	)
	if id, err = ctr.BaseController.GetID(c); err != nil {
		return account, err
	}
	if account, err = ctr.repo.GetByID(int(id)); err != nil {
		return account, err
	}

	return account, err
}

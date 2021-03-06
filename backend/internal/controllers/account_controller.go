package controllers

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/max-gryshin/local-chain/internal/contractions"
	"github.com/max-gryshin/local-chain/internal/dto"
	"github.com/max-gryshin/local-chain/internal/e"
	"github.com/max-gryshin/local-chain/internal/models"
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

// GetByID godoc
// @Summary      get account
// @Description  get account by id
// @Tags         account
// @Accept       json
// @Produce      json
// @Param		 id   path      int  true  "Account ID"
// @Success      200  {object}   dto.GetAccount
// @Security     ApiKeyAuth
// @Router       /api/account/{id} [get]
func (ctr *AccountController) GetByID(c echo.Context) error {
	var (
		err     error
		account models.Account
	)
	if account, err = ctr.getAccountByID(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.LoadGetAccountDTOFromModel(&account))
}

// GetAccounts   godoc
// @Summary      get accounts
// @Description  get accounts
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}    dto.GetAccounts
// @Security     ApiKeyAuth
// @Router       /api/account [get]
func (ctr *AccountController) GetAccounts(c echo.Context) error {
	var (
		accounts models.Accounts
		err      error
	)
	if accounts, err = ctr.repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadGetAccountsDTOCollectionFromModel(accounts))
}

// UpdateAccount godoc
// @Summary      update account
// @Description  update own account by id
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        message  body  dto.AccountOwnerUpdateRequest true  "AccountOwnerUpdateRequest"
// @Success      200  {object}  dto.AccountOwnerUpdate
// @Security     ApiKeyAuth
// @Router       /api/account/{id} [patch]
func (ctr *AccountController) UpdateAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
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

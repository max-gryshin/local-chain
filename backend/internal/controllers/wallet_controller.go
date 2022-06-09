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

// WalletController is HTTP controller for manage wallet
type WalletController struct {
	repo         contractions.WalletRepository
	errorHandler e.ErrorHandler
	BaseController
}

// NewWalletController return new instance of WalletController
func NewWalletController(repo contractions.WalletRepository, errorHandler e.ErrorHandler, v *validator.Validate) *WalletController {
	return &WalletController{
		repo:           repo,
		errorHandler:   errorHandler,
		BaseController: BaseController{*v},
	}
}

// GetByID       godoc
// @Summary      get wallet
// @Description  get wallet by id
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Wallet ID"
// @Success      200  {object}  dto.Wallet
// @Security     ApiKeyAuth
// @Router      /api/wallet/{id} [get]
func (ctr *WalletController) GetByID(c echo.Context) error {
	var (
		err    error
		wallet models.Wallet
	)
	if wallet, err = ctr.getWalletByID(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.LoadWalletDTOFromModel(&wallet))
}

// GetWallets    godoc
// @Summary      get wallets
// @Description  get wallets
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Success      200  {object} dto.Wallets
// @Security     ApiKeyAuth
// @Router       /api/wallet [get]
func (ctr *WalletController) GetWallets(c echo.Context) error {
	var (
		wallets models.Wallets
		err     error
	)
	if wallets, err = ctr.repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadWalletDTOCollectionFromModel(wallets))
}

func (ctr *WalletController) getWalletByID(c echo.Context) (models.Wallet, error) {
	var (
		id     int64
		err    error
		wallet models.Wallet
	)
	if id, err = ctr.BaseController.GetID(c); err != nil {
		return wallet, err
	}
	if wallet, err = ctr.repo.GetByID(int(id)); err != nil {
		return wallet, err
	}

	return wallet, err
}

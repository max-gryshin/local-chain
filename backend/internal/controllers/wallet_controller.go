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

// GetByID example: /api/wallet/{id}/
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

// GetWallets return list of users
// example: /api/wallet
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

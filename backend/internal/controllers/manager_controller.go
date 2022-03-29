package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/ZmaximillianZ/local-chain/internal/contractions"
	"github.com/ZmaximillianZ/local-chain/internal/dto"
	"github.com/ZmaximillianZ/local-chain/internal/e"
	"github.com/ZmaximillianZ/local-chain/internal/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// ManagerController is HTTP controller for manage user by manager
type ManagerController struct {
	repo         contractions.UserRepository
	errorHandler e.ErrorHandler
	BaseController
}

// NewManagerController return new instance of ManagerController
func NewManagerController(repo contractions.UserRepository, errorHandler e.ErrorHandler, v *validator.Validate) *ManagerController {
	return &ManagerController{
		repo:           repo,
		errorHandler:   errorHandler,
		BaseController: BaseController{*v},
	}
}

// Create        godoc
// @Summary      create user
// @Description  create user by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.User  true  "User"
// @Success      200  {object}  dto.User
// @Security     ApiKeyAuth
// @Router       /api/manager/user [post]
func (ctr *ManagerController) Create(c echo.Context) error {
	email := c.QueryParam("email")
	password := c.QueryParam("password")
	a := models.Auth{Email: email, Password: password}
	if errValidation := ctr.BaseController.validator.Struct(&a); errValidation != nil {
		return errValidation
	}
	var (
		userExist models.User
		err       error
		user      models.User
	)
	if userExist, err = ctr.repo.GetByEmail(a.Email); err != nil {
		return err
	}
	if userExist.ID != 0 {
		return errors.New("user with email " + a.Email + " exists")
	}
	user = models.User{Email: a.Email, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if errSetPassword := user.SetPassword(a.Password); errSetPassword != nil {
		return errSetPassword
	}
	if errCreateUser := ctr.repo.Create(&user); errCreateUser != nil {
		return errCreateUser
	}
	return c.JSON(http.StatusOK, dto.LoadUserDTOFromModel(&user))
}

// UpdateUser    godoc
// @Summary      update user
// @Description  update user by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.User  true  "User"
// @Success      200  {object}  dto.User
// @Security     ApiKeyAuth
// @Router       /api/manager/user/{id} [patch]
func (ctr *UserController) UpdateUser(c echo.Context) error {
	var (
		err  error
		user models.User
	)
	if user, err = ctr.getUserByID(c); err != nil {
		return err
	}
	dtoUser := dto.LoadUserDTOFromModel(&user)
	if errBindOrValidate := ctr.BindAndValidate(c, dtoUser); errBindOrValidate != nil {
		return errBindOrValidate
	}
	if errUpdateUser := ctr.repo.Update(dto.LoadUserModelFromDTO(dtoUser)); errUpdateUser != nil {
		return errUpdateUser
	}
	return c.JSON(http.StatusOK, dtoUser)
}

// GetOrdersByManager godoc
// @Summary           get orders
// @Description       get orders by manager
// @Tags              order
// @Accept            json
// @Produce           json
// @Success           200  {object} dto.Orders
// @Security          ApiKeyAuth
// @Router            /api/manager/order [get]
func (ctr *OrderController) GetOrdersByManager(c echo.Context) error {
	var (
		orders models.Orders // todo: change to orders
		err    error
	)
	if orders, err = ctr.repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadOrderDTOCollectionFromModel(orders))
}

// HandleOrder   godoc
// @Summary      handle order
// @Description  handle order by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Order  true  "Order"
// @Success      200  {object} dto.Order
// @Security     ApiKeyAuth
// @Router       /api/manager/order/{orderId} [patch]
func (ctr *ManagerController) HandleOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// CreateAccount godoc
// @Summary      create an account
// @Description  creating an account by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Account  true  "Account"
// @Success      200  {object} dto.Account
// @Security     ApiKeyAuth
// @Router       /api/manager/account/{userid} [post]
func (ctr *ManagerController) CreateAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// UpdateAccount godoc
// @Summary      update an account
// @Description  updating an account by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Account  true  "Account"
// @Success      200  {object} dto.Account
// @Security     ApiKeyAuth
// @Router       /api/manager/account/{accountId} [patch]
func (ctr *ManagerController) UpdateAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// CreateWallet  godoc
// @Summary      create a wallet
// @Description  creating a wallet by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Wallet  true  "Wallet"
// @Success      200  {object} dto.Wallet
// @Security     ApiKeyAuth
// @Router       /api/manager/wallet [post]
func (ctr *ManagerController) CreateWallet(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// UpdateWallet  godoc
// @Summary      update a wallet
// @Description  updating a wallet by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Wallet  true  "Wallet"
// @Success      200  {object} dto.Wallet
// @Security     ApiKeyAuth
// @Router       /api/manager/wallet/{walletId} [patch]
func (ctr *ManagerController) UpdateWallet(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// Debit         godoc
// @Summary      debit
// @Description  debit amount from user wallet
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Wallet  true  "Wallet"
// @Success      200  {object} dto.Wallet
// @Security     ApiKeyAuth
// @Router       /api/manager/wallet/{walletId}/debit [post]
// todo: create dto
func (ctr *ManagerController) Debit(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// Credit        godoc
// @Summary      credit
// @Description  credit amount from user wallet
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Wallet  true  "Wallet"
// @Success      200  {object} dto.Wallet
// @Security     ApiKeyAuth
// @Router       /api/manager/wallet/{walletId}/credit [post]
// todo: create dto
func (ctr *ManagerController) Credit(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/max-gryshin/local-chain/internal/contractions"
	"github.com/max-gryshin/local-chain/internal/dto"
	"github.com/max-gryshin/local-chain/internal/e"
	"github.com/max-gryshin/local-chain/internal/middleware/access"
	"github.com/max-gryshin/local-chain/internal/models"
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
// @Param        message  body  dto.UserRegistration  true  "User"
// @Success      200  {object}  dto.UserRegistration
// @Security     ApiKeyAuth
// @Router       /api/manager/user [post]
func (ctr *ManagerController) Create(c echo.Context) error {
	var (
		managerID int
		err       error
		userExist models.User
	)
	managerIDString := c.Get(access.UserID).(string)
	if managerID, err = strconv.Atoi(managerIDString); err != nil {
		return err
	}
	newUserDTO := dto.UserRegistration{}
	if errBinding := c.Bind(&newUserDTO); errBinding != nil {
		return errBinding
	}
	if errValidate := ctr.validator.Struct(newUserDTO); errValidate != nil {
		return errValidate
	}
	if userExist, err = ctr.repo.GetByEmail(*newUserDTO.Email); err != nil {
		return err
	}
	if userExist.ID != 0 {
		return errors.New("user with email " + *newUserDTO.Email + " exists")
	}
	user := dto.LoadUserModelFromUserRegistrationDTO(&newUserDTO)
	user.CreatedBy = managerID
	user.UpdatedBy = managerID
	user.ManagerID = managerID
	if errSetPassword := user.SetPassword(newUserDTO.Password); errSetPassword != nil {
		return errSetPassword
	}
	if errCreateUser := ctr.repo.Create(user); errCreateUser != nil {
		return errCreateUser
	}
	return c.JSON(http.StatusOK, dto.LoadUserByManagerDTOFromModel(user))
}

// UpdateUser    godoc
// @Summary      update user
// @Description  update user by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "User ID"
// @Param        message  body  dto.UserByManager  true  "User"
// @Success      200  {object}  dto.UserByManager
// @Security     ApiKeyAuth
// @Router       /api/manager/user/{id} [patch]
func (ctr *UserController) UpdateUser(c echo.Context) error {
	var (
		managerID int
		err       error
		user      models.User
	)
	if user, err = ctr.getUserByID(c); err != nil {
		return err
	}
	managerIDString := c.Get(access.UserID).(string)
	if managerID, err = strconv.Atoi(managerIDString); err != nil {
		return err
	}
	dtoUser := dto.LoadUserByManagerDTOFromModel(&user)
	if errBindOrValidate := ctr.BindAndValidate(c, dtoUser); errBindOrValidate != nil {
		return errBindOrValidate
	}
	dtoUser.UpdatedBy = managerID
	if errUpdateUser := ctr.Repo.Update(dto.LoadUserModelFromUserByManagerDTO(dtoUser)); errUpdateUser != nil {
		return errUpdateUser
	}
	return c.JSON(http.StatusOK, dtoUser)
}

// GetUsersByManager godoc
// @Summary          get users by manager
// @Description      get users by manager
// @Tags             manager
// @Accept           json
// @Produce          json
// @Success          200  {object} dto.UsersByManager
// @Security         ApiKeyAuth
// @Router           /api/manager/user [get]
func (ctr *UserController) GetUsersByManager(c echo.Context) error {
	var (
		users models.Users
		err   error
	)
	if users, err = ctr.Repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadUsersByManagerDTOCollectionFromModel(users))
}

// GetUserByID   godoc
// @Summary      get user
// @Description  get user by id
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.UserByManager
// @Security     ApiKeyAuth
// @Router       /api/manager/user/{id} [get]
func (ctr *UserController) GetUserByID(c echo.Context) error {
	var (
		err  error
		user models.User
	)
	if user, err = ctr.getUserByID(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.LoadUserByManagerDTOFromModel(&user))
}

// GetOrdersByManager godoc
// @Summary           get orders
// @Description       get orders by manager
// @Tags              order
// @Accept            json
// @Produce           json
// @Success           200  {object} dto.GetOrders
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

	return c.JSON(http.StatusOK, dto.LoadGetOrdersDTOCollectionFromModel(orders))
}

// HandleOrder   godoc
// @Summary      handle order
// @Description  handle order by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.OrderByManager  true  "GetOrder"
// @Success      200  {object} dto.GetOrder
// @Security     ApiKeyAuth
// @Router       /api/manager/order/{id} [patch]
func (ctr *ManagerController) HandleOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// CreateAccount godoc
// @Summary      create an account
// @Description  creating an account by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.AccountByManager  true  "AccountByManager"
// @Success      200  {object} dto.GetAccount
// @Security     ApiKeyAuth
// @Router       /api/manager/account [post]
func (ctr *ManagerController) CreateAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// UpdateAccount godoc
// @Summary      update an account
// @Description  updating an account by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.AccountByManager  true  "AccountByManager"
// @Success      200  {object} dto.GetAccount
// @Security     ApiKeyAuth
// @Router       /api/manager/account/{id} [patch]
func (ctr *ManagerController) UpdateAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// CreateWallet  godoc
// @Summary      create a wallet
// @Description  creating a wallet by manager
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.WalletCreate  true  "Wallet"
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
// @Param        message  body  dto.WalletCreate  true  "Wallet"
// @Success      200  {object} dto.Wallet
// @Security     ApiKeyAuth
// @Router       /api/manager/wallet/{id} [patch]
func (ctr *ManagerController) UpdateWallet(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// Debit         godoc
// @Summary      debit
// @Description  debit amount from user wallet
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Debit  true  "Debit"
// @Success      200  {object} dto.Wallet
// @Security     ApiKeyAuth
// @Router       /api/manager/wallet/{id}/debit [post]
func (ctr *ManagerController) Debit(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// Credit        godoc
// @Summary      credit
// @Description  credit amount from user wallet
// @Tags         manager
// @Accept       json
// @Produce      json
// @Param        message  body  dto.Credit  true  "Credit"
// @Success      200  {object} dto.Wallet
// @Security     ApiKeyAuth
// @Router       /api/manager/wallet/{id}/credit [post]
func (ctr *ManagerController) Credit(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

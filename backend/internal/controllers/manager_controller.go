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

// Create create user
// example: /api/manager/user
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

// GetOrders example: /api/manager/order
func (ctr *ManagerController) GetOrders(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// GetOrder example: /api/manager/order/{orderId}
func (ctr *ManagerController) GetOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// HandleOrder example: /api/manager/order/{orderId}
func (ctr *ManagerController) HandleOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// Debit example: /api/manager/wallet/{walletId}/debit
func (ctr *ManagerController) Debit(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

// Credit example: /api/manager/wallet/{walletId}/credit
func (ctr *ManagerController) Credit(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"1", "two", "110"})
}

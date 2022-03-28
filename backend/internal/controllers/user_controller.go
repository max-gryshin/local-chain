package controllers

import (
	"errors"

	"github.com/ZmaximillianZ/local-chain/internal/dto"
	"github.com/go-playground/validator"

	"github.com/ZmaximillianZ/local-chain/internal/e"
	"github.com/ZmaximillianZ/local-chain/internal/models"
	"github.com/ZmaximillianZ/local-chain/internal/utils"
	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/ZmaximillianZ/local-chain/internal/contractions"
)

// UserController is HTTP controller for manage users
type UserController struct {
	repo         contractions.UserRepository
	errorHandler e.ErrorHandler
	BaseController
}

// NewUserController return new instance of UserController
func NewUserController(repo contractions.UserRepository, errorHandler e.ErrorHandler, v *validator.Validate) *UserController {
	return &UserController{
		repo:           repo,
		errorHandler:   errorHandler,
		BaseController: BaseController{*v},
	}
}

// GetByID return user by id
// example: /api/user/{id}/
func (ctr *UserController) GetByID(c echo.Context) error {
	var (
		err  error
		user models.User
	)
	if user, err = ctr.getUserByID(c); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.LoadUserDTOFromModel(&user)) // todo: is it have sense?
}

// Authenticate @Summary Authenticate
// description: user authorization
// example: /api/auth
func (ctr *UserController) Authenticate(c echo.Context) error {
	email := c.QueryParam("email")
	password := c.QueryParam("password")
	a := models.Auth{Email: email, Password: password}
	if errValidation := ctr.BaseController.validator.Struct(&a); errValidation != nil {
		return errValidation
	}
	var (
		user  models.User
		err   error
		token string
	)
	if user, err = ctr.repo.GetByEmail(email); err != nil {
		return err
	}
	if user.InvalidPassword(password) {
		return errors.New("invalid password")
	}
	if token, err = utils.GenerateToken(a.Email, a.Password, user.ID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// GetUsers return list of users
// example: /api/user/all
func (ctr *UserController) GetUsers(c echo.Context) error {
	var (
		users models.Users
		err   error
	)
	if users, err = ctr.repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadUserDTOCollectionFromModel(users))
}

func (ctr *UserController) getUserByID(c echo.Context) (models.User, error) {
	var (
		id   int64
		err  error
		user models.User
	)
	if id, err = ctr.BaseController.GetID(c); err != nil {
		return user, err
	}
	if user, err = ctr.repo.GetByID(int(id)); err != nil {
		return user, err
	}

	return user, err
}

// Update return user by id
// example: /api/manager/user/{id}/
func (ctr *UserController) Update(c echo.Context) error {
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

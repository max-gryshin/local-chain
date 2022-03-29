package controllers

import (
	"errors"

	"github.com/ZmaximillianZ/local-chain/internal/dto"
	"github.com/ZmaximillianZ/local-chain/internal/utils"
	"github.com/go-playground/validator"

	"github.com/ZmaximillianZ/local-chain/internal/e"
	"github.com/ZmaximillianZ/local-chain/internal/models"
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

// Authenticate  godoc
// @Summary      authenticate
// @Description  authenticate user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        email     path  string  true  "email"
// @Param        password  path  string  true  "password"
// @Success      200  {object}  dto.User
// @Router       /api/auth [post]
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

// GetByID godoc
// @Summary      get user
// @Description  get user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.User
// @Security     ApiKeyAuth
// @Router       /api/user/{id} [get]
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

// GetUsers      godoc
// @Summary      get users
// @Description  get users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object} dto.Users
// @Security     ApiKeyAuth
// @Router       /api/user [get]
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

// Update godoc
// @Summary      update user
// @Description  update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        message  body  dto.User  true  "User"
// @Success      200  {object}  dto.User
// @Security     ApiKeyAuth
// @Router       /api/user/{id} [patch]
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

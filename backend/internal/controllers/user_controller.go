package controllers

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/max-gryshin/local-chain/internal/dto"
	"github.com/max-gryshin/local-chain/internal/middleware/access"
	"github.com/max-gryshin/local-chain/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/max-gryshin/local-chain/internal/e"
	"github.com/max-gryshin/local-chain/internal/models"

	"net/http"

	"github.com/max-gryshin/local-chain/internal/contractions"
)

// UserController is HTTP controller for manage users
type UserController struct {
	Repo         contractions.UserRepository
	errorHandler e.ErrorHandler
	BaseController
}

// NewUserController return new instance of UserController
func NewUserController(repo contractions.UserRepository, errorHandler e.ErrorHandler, v *validator.Validate) *UserController {
	return &UserController{
		Repo:           repo,
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
// @Param        email     query  string  true  "email"
// @Param        password  query  string  true  "password"
// @Success      200  {string}  dto.Authenticate
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
	if user, err = ctr.Repo.GetByEmail(email); err != nil {
		return err
	}
	if user.InvalidPassword(password) {
		return errors.New("invalid password")
	}
	if token, err = utils.GenerateToken(user.ID); err != nil {
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
// @Success      200  {object}  dto.GetUserOwner
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
	return c.JSON(http.StatusOK, dto.LoadGetUserOwnerDTOFromModel(&user))
}

// GetUsers      godoc
// @Summary      get users
// @Description  get users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object} dto.GetUserOwners
// @Security     ApiKeyAuth
// @Router       /api/user [get]
func (ctr *UserController) GetUsers(c echo.Context) error {
	var (
		users models.Users
		err   error
	)
	if users, err = ctr.Repo.GetAll(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoadGetUserOwnerDTOCollectionFromModel(users))
}

// Update godoc
// @Summary      update user
// @Description  update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        message  body  dto.UpdateUserOwnerRequest  true  "User"
// @Success      200  {object}  dto.UpdateUserOwnerRequest
// @Security     ApiKeyAuth
// @Router       /api/user [patch]
func (ctr *UserController) Update(c echo.Context) error {
	var (
		userID int
		err    error
		user   models.User
	)
	id := c.Get(access.UserID).(string)
	if userID, err = strconv.Atoi(id); err != nil {
		return err
	}
	if user, err = ctr.Repo.GetByID(userID); err != nil {
		return err
	}
	dtoUser := dto.LoadUpdateUserOwnerDTOFromModel(&user)
	if errBindOrValidate := ctr.BindAndValidate(c, dtoUser); errBindOrValidate != nil {
		return errBindOrValidate
	}
	newModel := dto.LoadUserModelFromUpdateUserOwnerDTO(dtoUser)
	newModel.UpdatedBy = userID
	newModel.UpdatedAt = time.Now()
	if errUpdateUser := ctr.Repo.Update(newModel); errUpdateUser != nil {
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
	if user, err = ctr.Repo.GetByID(int(id)); err != nil {
		return user, err
	}

	return user, err
}

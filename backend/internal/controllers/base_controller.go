package controllers

import (
	"errors"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-playground/validator"

	"github.com/labstack/echo/v4"
)

type BaseController struct {
	validator validator.Validate
}

func (ctr *BaseController) GetID(c echo.Context) (int64, error) {
	var (
		id  int64
		err error
	)
	if id, err = strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		return id, errors.New("invalid id")
	}

	return id, nil
}

func (ctr *BaseController) BindAndValidate(c echo.Context, model interface{}) error {
	testM := model
	if errBinding := c.Bind(&model); errBinding != nil {
		return errBinding
	}
	if errValidate := ctr.validator.Struct(testM); errValidate != nil { //todo: fix for workout
		return errValidate
	}

	return nil
}

func (ctr *BaseController) GetUserIDFromToken(c echo.Context) (int, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	idFromClaim := claims["id"].(string)
	var (
		ID  int
		err error
	)
	if ID, err = strconv.Atoi(idFromClaim); err != nil {
		return ID, err
	}

	return ID, nil
}

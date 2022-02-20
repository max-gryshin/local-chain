package jwt

import (
	"net/http"

	"github.com/ZmaximillianZ/local-chain/internal/e"
	"github.com/ZmaximillianZ/local-chain/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// JWT is jwt middleware
func JWT(c echo.Context) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		var code int
		var data interface{}
		code = e.Success
		// TODO: recommended use Authentication: Bearer
		// https://swagger.io/docs/specification/authentication/bearer-authentication/
		token := c.Request().Header.Get("X-AUTH-TOKEN")
		if token == "" {
			code = e.InvalidParams
		} else {
			_, err := utils.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ErrorAuthCheckTokenTimeout
				default:
					code = e.ErrorAuthCheckTokenFail
				}
			}
		}

		if code != e.Success {
			return func(c echo.Context) error {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"code": code,
					"msg":  e.GetMsg(code),
					"data": data,
				})
			}
		}

		return h
	}
}

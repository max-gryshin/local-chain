package access

import (
	"strings"

	"github.com/ZmaximillianZ/local-chain/internal/utils"
	"github.com/labstack/echo/v4"
)

const UserID = "userID"

// IsResourceAvailable middleware check is caller have access to this request
func IsResourceAvailable(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqToken := c.Request().Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		userID, err := utils.GetAuthenticatedUserID(splitToken[1])
		if err != nil {
			return err
		}
		paramID := c.Param("userId")
		if paramID != "" && c.Param("id") != userID {
			// todo: check caller is a manager and have a access to resource
			return echo.ErrForbidden
		}
		c.Set("userID", userID)

		return next(c)
	}
}

package access

import (
	"strconv"
	"strings"

	"github.com/ZmaximillianZ/local-chain/internal/contractions"
	"github.com/ZmaximillianZ/local-chain/internal/models"

	"github.com/ZmaximillianZ/local-chain/internal/utils"
	"github.com/labstack/echo/v4"
)

const UserID = "userID"

type ResourceAccess struct {
	repo contractions.UserRepository
}

func NewResourceAccess(repo contractions.UserRepository) *ResourceAccess {
	return &ResourceAccess{repo: repo}
}

// IsResourceAvailable middleware check is caller have access to this request
func (ra *ResourceAccess) IsResourceAvailable(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqToken := c.Request().Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		userID, err := utils.GetAuthenticatedUserID(splitToken[1])
		if err != nil {
			return err
		}
		paramID := c.Param("id")
		if paramID != "" && paramID != userID {
			isManager, err := ra.isManager(paramID, userID)
			if err != nil {
				return err
			}
			if !isManager {
				return echo.ErrForbidden
			}
		}
		c.Set("userID", userID)

		return next(c)
	}
}

func (ra *ResourceAccess) isManager(requiredID, authenticatedUserID string) (bool, error) {
	authenticatedUserIDint, err := strconv.Atoi(authenticatedUserID)
	if err != nil {
		return false, err
	}
	isManager, err := ra.repo.GetManagerIDs()
	if err != nil {
		return false, err
	}
	if isManager == nil {
		return false, nil
	}
	if !utils.ContainsInt(isManager, authenticatedUserIDint) {
		return false, nil
	}
	userIDint, err := strconv.Atoi(requiredID)
	if err != nil {
		return false, err
	}
	user, err := ra.repo.GetByID(userIDint)
	if err != nil {
		return false, err
	}
	authenticatedUser, err := ra.repo.GetByID(authenticatedUserIDint)
	if err != nil {
		return false, err
	}
	if user.ManagerID != authenticatedUser.ID {
		// it means we must have user with admin role
		if utils.ContainsString(user.Roles, models.RoleAdmin) {
			return false, nil
		}
		return ra.isManager(strconv.Itoa(user.ManagerID), authenticatedUserID)
	}
	return true, nil
}

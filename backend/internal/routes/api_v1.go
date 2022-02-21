package routes

import (
	"github.com/ZmaximillianZ/local-chain/internal/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterAPIV1 initialize routing information
func RegisterAPIV1(
	router *echo.Group,
	userController *controllers.UserController,
) {
	jwt := middleware.JWT([]byte("get_key_from_env"))

	router.POST("/auth", userController.Authenticate)
	router.POST("/create", userController.Create)
	user := router.Group("/users")
	user.GET("/", userController.GetUsers, jwt)
	user.POST("/", userController.Create)
	user.GET("/:id", userController.GetByID, jwt)
	user.PUT("/:id", userController.Update, jwt)
	user.DELETE("/:id", userController.Delete, jwt)
}

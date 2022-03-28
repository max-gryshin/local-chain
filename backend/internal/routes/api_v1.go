package routes

import (
	"github.com/ZmaximillianZ/local-chain/internal/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterAPI initialize routing information
func RegisterAPI(
	router *echo.Group,
	userController *controllers.UserController,
	managerController *controllers.ManagerController,
	accountController *controllers.AccountController,
	orderController *controllers.OrderController,
	walletController *controllers.WalletController,
	transactionController *controllers.TransactionController,
) {
	jwt := middleware.JWT([]byte("get_key_from_env"))

	router.POST("/auth", userController.Authenticate)
	user := router.Group("/user")
	user.GET("/all", userController.GetUsers, jwt)
	user.GET("/:id", userController.GetByID, jwt)

	account := router.Group("/account")
	account.GET("/all", accountController.GetAccounts, jwt)
	account.GET("/:id", accountController.GetByID, jwt)
	account.GET("/:id/cash-out", accountController.CashOut, jwt)

	order := router.Group("/order")
	order.GET("/all", orderController.GetOrders, jwt)
	order.GET("/:id", orderController.GetByID, jwt)

	wallet := router.Group("/wallet")
	wallet.GET("/:id", walletController.GetByID, jwt)
	wallet.GET("/all", walletController.GetWallets, jwt)

	transaction := router.Group("/transaction")
	transaction.POST("/", transactionController.SendTransaction, jwt)
	transaction.GET("/all", transactionController.GetTransactions, jwt)
	transaction.GET("/:userId/all", transactionController.GetUserTransactions, jwt)

	manager := router.Group("/manager")
	manager.POST("/", managerController.Create, jwt)                        // isManagerMiddleware
	manager.PATCH("/:id", userController.Update, jwt)                       // isManagerMiddleware
	manager.GET("/order", managerController.GetOrders)                      // isManagerMiddleware
	manager.GET("/order/:orderId", managerController.GetOrders, jwt)        // isManagerMiddleware
	manager.POST("/order/:orderId", managerController.HandleOrder, jwt)     // isManagerMiddleware
	manager.POST("/wallet/debit/:walletId", managerController.Debit, jwt)   // isManagerMiddleware
	manager.POST("/wallet/credit/:walletId", managerController.Credit, jwt) // isManagerMiddleware
}

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

	user := router.Group("/user", jwt)
	user.GET("/:id", userController.GetByID)
	user.GET("", userController.GetUsers)
	user.PATCH("", userController.Update)

	account := router.Group("/account", jwt)
	account.GET("/:accountId", accountController.GetByID)
	account.GET("", accountController.GetAccounts)
	account.PATCH("/:accountId", accountController.UpdateAccount)
	account.GET("/:accountId/cash-out", accountController.CashOut)

	order := router.Group("/order", jwt)
	order.GET("/:orderId", orderController.GetByID)
	order.GET("", orderController.GetOrders)

	wallet := router.Group("/wallet", jwt)
	wallet.GET("/:walletId", walletController.GetByID)
	wallet.GET("", walletController.GetWallets)

	transaction := router.Group("/transaction", jwt)

	transaction.GET("", transactionController.GetTransactions)
	transaction.GET("/:userId", transactionController.GetUserTransactions)
	transaction.POST("", transactionController.SendTransaction)

	manager := router.Group("/manager", jwt)
	manager.POST("manager/user", managerController.Create)
	manager.PATCH("/user/:userId", userController.Update)
	manager.GET("/order", orderController.GetOrdersByManager)
	manager.PATCH("/order/:orderId", managerController.HandleOrder)
	manager.POST("/account/:userId", managerController.CreateAccount)
	manager.PATCH("/account/:accountId", managerController.UpdateAccount)
	manager.POST("/wallet", managerController.CreateWallet)
	manager.PATCH("/wallet/:walletId", managerController.UpdateWallet)
	manager.POST("/wallet/:walletId/debit", managerController.Debit)
	manager.POST("/wallet/:walletId/credit", managerController.Credit)
}

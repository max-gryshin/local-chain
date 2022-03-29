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
	user.GET("/:id", userController.GetByID, jwt)
	user.GET("", userController.GetUsers, jwt)
	user.PATCH("/:id", userController.Update, jwt)

	account := router.Group("/account")
	account.GET("/:id", accountController.GetByID, jwt)
	account.GET("", accountController.GetAccounts, jwt)
	account.PATCH("/:id", accountController.UpdateAccount, jwt)
	account.GET("/:id/cash-out", accountController.CashOut, jwt)

	order := router.Group("/order")
	order.GET("/:id", orderController.GetByID, jwt)
	order.GET("", orderController.GetOrders, jwt)

	wallet := router.Group("/wallet")
	wallet.GET("/:id", walletController.GetByID, jwt)
	wallet.GET("", walletController.GetWallets, jwt)

	transaction := router.Group("/transaction")

	transaction.GET("", transactionController.GetTransactions, jwt)
	transaction.GET("/:userId", transactionController.GetUserTransactions, jwt)
	transaction.POST("", transactionController.SendTransaction, jwt)

	manager := router.Group("/manager")
	manager.POST("/user", managerController.Create)                            // isManagerMiddleware
	manager.PATCH("/user/:id", userController.Update, jwt)                     // isManagerMiddleware
	manager.GET("/order", orderController.GetOrdersByManager, jwt)             // isManagerMiddleware
	manager.PATCH("/order/:orderId", managerController.HandleOrder, jwt)       // isManagerMiddleware
	manager.POST("/account/:userId", managerController.CreateAccount, jwt)     // isManagerMiddleware
	manager.PATCH("/account/:accountId", managerController.UpdateAccount, jwt) // isManagerMiddleware
	manager.POST("/wallet", managerController.CreateWallet, jwt)               // isManagerMiddleware
	manager.PATCH("/wallet/:walletId", managerController.UpdateWallet, jwt)    // isManagerMiddleware
	manager.POST("/wallet/:walletId/debit", managerController.Debit, jwt)      // isManagerMiddleware
	manager.POST("/wallet/:walletId/credit", managerController.Credit, jwt)    // isManagerMiddleware
}

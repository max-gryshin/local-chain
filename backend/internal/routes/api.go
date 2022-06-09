package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/max-gryshin/local-chain/internal/controllers"
	"github.com/max-gryshin/local-chain/internal/middleware/access"
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

	resourceAccess := access.NewResourceAccess(userController.Repo)
	user := router.Group("/user", jwt)
	user.GET("/:id", userController.GetByID)
	user.GET("", userController.GetUsers)
	user.PATCH("", userController.Update, resourceAccess.IsResourceAvailable)

	account := router.Group("/account", jwt)
	account.GET("/:id", accountController.GetByID)
	account.GET("", accountController.GetAccounts)
	account.PATCH("/:id", accountController.UpdateAccount)

	order := router.Group("/order", jwt)
	order.GET("/:id", orderController.GetByID)
	order.GET("", orderController.GetOrders)
	order.POST("/:id/cash-out", orderController.CashOut)

	wallet := router.Group("/wallet", jwt)
	wallet.GET("/:id", walletController.GetByID)
	wallet.GET("", walletController.GetWallets)

	transaction := router.Group("/transaction", jwt)

	transaction.GET("", transactionController.GetTransactions)
	transaction.GET("/user/:id", transactionController.GetUserTransactions)
	transaction.POST("", transactionController.SendTransaction)

	manager := router.Group("/manager", jwt)
	manager.POST("/user", managerController.Create, resourceAccess.IsResourceAvailable)
	manager.PATCH("/user/:id", userController.UpdateUser, resourceAccess.IsResourceAvailable)
	manager.GET("/user", userController.GetUsersByManager)
	manager.GET("/user/:id", userController.GetUserByID)
	manager.GET("/order", orderController.GetOrdersByManager)
	manager.PATCH("/order/:id", managerController.HandleOrder)
	manager.POST("/account", managerController.CreateAccount)
	manager.PATCH("/account/:id", managerController.UpdateAccount)
	manager.POST("/wallet", managerController.CreateWallet)
	manager.PATCH("/wallet/:id", managerController.UpdateWallet)
	manager.POST("/wallet/:id/debit", managerController.Debit)
	manager.POST("/wallet/:id/credit", managerController.Credit)
}

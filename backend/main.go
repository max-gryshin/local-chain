package main

import (
	"log"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/max-gryshin/local-chain/docs"
	"github.com/max-gryshin/local-chain/internal/controllers"
	"github.com/max-gryshin/local-chain/internal/db"
	errorhandler "github.com/max-gryshin/local-chain/internal/e"
	"github.com/max-gryshin/local-chain/internal/logging"
	"github.com/max-gryshin/local-chain/internal/repository"
	"github.com/max-gryshin/local-chain/internal/routes"
	"github.com/max-gryshin/local-chain/internal/setting"
)

// @title        Local chain API
// @version      1.0
// @description  This is a local chain server.

// @contact.name   Maxim Hryshyn
// @contact.email  goooglemax1993@gmail.com
// @schemes http
// @host      	   0.0.0.0:1323
// @BasePath       /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(err.Error())
		log.Print("No .env file found")
	}
	settings := setting.LoadSetting()

	logging.Setup(&settings.App)
	defer func() {
		err := logging.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	dbContext, err := db.CreateDatabaseContext(settings.DBConfig)
	if err != nil {
		log.Print(err)
		panic(err)
	}
	userRepo := repository.NewUserRepository(dbContext.Connection, dbContext.QueryBuilder)
	accountRepo := repository.NewAccountRepository(dbContext.Connection, dbContext.QueryBuilder)
	orderRepo := repository.NewOrderRepository(dbContext.Connection, dbContext.QueryBuilder)
	walletRepo := repository.NewWalletRepository(dbContext.Connection, dbContext.QueryBuilder)
	eHandler := errorhandler.ErrorHandler{}
	validate := validator.New()
	userController := controllers.NewUserController(userRepo, eHandler, validate)
	managerController := controllers.NewManagerController(userRepo, eHandler, validate)
	accountController := controllers.NewAccountController(accountRepo, eHandler, validate)
	orderController := controllers.NewOrderController(orderRepo, eHandler, validate)
	walletController := controllers.NewWalletController(walletRepo, eHandler, validate)
	transactionController := controllers.NewTransactionController(eHandler, validate)
	eco := echo.New()
	eco.Debug = true
	eco.GET("/swagger/*", echoSwagger.WrapHandler)
	router := eco.Group("/api")
	router.Use(eHandler.Handle)
	routes.RegisterAPI(
		router,
		userController,
		managerController,
		accountController,
		orderController,
		walletController,
		transactionController,
	)
	// todo: use with TLS StartTLS(":1323", "cert.pem", "key.pem")
	eco.Logger.Fatal(eco.Start(":1323"))
}

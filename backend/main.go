package main

import (
	"log"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/ZmaximillianZ/local-chain/docs"
	"github.com/ZmaximillianZ/local-chain/internal/controllers"
	"github.com/ZmaximillianZ/local-chain/internal/db"
	errorhandler "github.com/ZmaximillianZ/local-chain/internal/e"
	"github.com/ZmaximillianZ/local-chain/internal/logging"
	"github.com/ZmaximillianZ/local-chain/internal/repository"
	"github.com/ZmaximillianZ/local-chain/internal/routes"
	"github.com/ZmaximillianZ/local-chain/internal/setting"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
	userController := controllers.NewUserController(userRepo, eHandler, validator.New())
	managerController := controllers.NewManagerController(userRepo, eHandler, validator.New())
	accountController := controllers.NewAccountController(accountRepo, eHandler, validator.New())
	orderController := controllers.NewOrderController(orderRepo, eHandler, validator.New())
	walletController := controllers.NewWalletController(walletRepo, eHandler, validator.New())
	transactionController := controllers.NewTransactionController(eHandler, validator.New())
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

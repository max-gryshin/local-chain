package main

import (
	"log"

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
	eHandler := errorhandler.ErrorHandler{}
	userController := controllers.NewUserController(userRepo, eHandler, validator.New())
	eco := echo.New()
	router := eco.Group("/api/v1")
	router.Use(eHandler.Handle)
	routes.RegisterAPIV1(router, userController)
	// todo: use with TLS StartTLS(":1323", "cert.pem", "key.pem")
	eco.Logger.Fatal(eco.Start(":1323"))
}

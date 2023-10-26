package main

import (
	"fmt"
	"miniproject-alterra/app"
	"miniproject-alterra/app/config"
	"miniproject-alterra/app/database"
	app_validator "miniproject-alterra/app/validator"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	cfg := config.GetConfig()
	db := database.InitDBMysql(cfg)

	echo := echo.New()
	echo.Validator = &app_validator.CustomValidator{Validator: validator.New()}
	echo.Pre(middleware.RemoveTrailingSlash())
	echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] status=${status} method=${method} uri=${uri} latency=${latency_human} \n",
	}))

	v1 := echo.Group("/v1")
	app.Bootstrap(db, v1, cfg)

	host := os.Getenv("HOST")
	if host == "" {
		host = "8080"
	}
	echo.Logger.Fatal(echo.Start(fmt.Sprintf(":%s", host)))

}

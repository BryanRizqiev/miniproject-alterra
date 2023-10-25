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

	e := echo.New()
	e.Validator = &app_validator.CustomValidator{Validator: validator.New()}
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] status=${status} method=${method} uri=${uri} latency=${latency_human} \n",
	}))

	app.Bootstrap(db, e, cfg)

	host := os.Getenv("HOST")
	if host == "" {
		host = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", host)))

}

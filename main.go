package main

import (
	"fmt"
	"miniproject-alterra/app"
	"miniproject-alterra/app/config"
	"miniproject-alterra/app/database"
	"miniproject-alterra/app/lib"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func main() {

	cfg := config.GetConfig()
	db := database.InitDBMysql(cfg)
	// if err := database.InitMigrationMySQL(db); err != nil {
	// 	panic(err)
	// }

	e := echo.New()
	e.Validator = &lib.CustomValidator{Validator: validator.New()}

	app.Bootstrap(db, e)

	host := os.Getenv("HOST")
	if host == "" {
		host = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", host)))

}

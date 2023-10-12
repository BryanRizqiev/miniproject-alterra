package app

import (
	user_controller "miniproject-alterra/module/user/controller"
	mysql_user_repository "miniproject-alterra/module/user/repository/mysql"
	user_service "miniproject-alterra/module/user/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Bootstrap(db *gorm.DB, e *echo.Echo) {

	userRepository := mysql_user_repository.NewUserRepository(db)
	userService := user_service.NewUserService(userRepository)
	userController := user_controller.NewUserController(userService)

	e.POST("/register", userController.Register)

}

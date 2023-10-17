package app

import (
	"miniproject-alterra/app/config"
	global_service "miniproject-alterra/module/global/service"
	user_controller "miniproject-alterra/module/user/controller"
	mysql_user_repository "miniproject-alterra/module/user/repository/mysql"
	user_service "miniproject-alterra/module/user/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Bootstrap(db *gorm.DB, e *echo.Echo, config *config.AppConfig) {

	userRepository := mysql_user_repository.NewUserRepository(db)
	userService := user_service.NewUserService(userRepository)

	emailService := global_service.NewEmailService(config)
	userController := user_controller.NewUserController(userService, emailService)

	e.POST("/register", userController.Register)

}

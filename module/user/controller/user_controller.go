package user_controller

import (
	global_entity "miniproject-alterra/module/global/entity"
	user_request "miniproject-alterra/module/user/controller/request"
	user_response "miniproject-alterra/module/user/controller/response"
	user_entity "miniproject-alterra/module/user/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService  user_entity.UserServiceInterface
	emailService global_entity.EmailServiceInterface
	storageService global_entity.StorageServiceInterface
}

func NewUserController(userService user_entity.UserServiceInterface, emailService global_entity.EmailServiceInterface, storageService global_entity.StorageServiceInterface) *UserController {

	return &UserController{
		userService:  userService,
		emailService: emailService,
		storageService: storageService
	}

}

func (this *UserController) Register(ctx echo.Context) error {

	req := new(user_request.RegisterRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, user_response.RegisterResponse{
			Message: "Request not valid",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	userDTO := user_entity.UserDTO{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Address:  req.Address,
	}

	if err := this.userService.Register(userDTO); err != nil {
		return ctx.JSON(http.StatusInternalServerError, user_response.RegisterResponse{
			Message: "Server error",
		})
	}

	format := global_entity.SendEmailFormat{
		To:      userDTO.Email,
		Cc:      userDTO.Email,
		Subject: "Test",
		Body:    "Registration Success",
	}
	go this.emailService.SendEmail(format)

	return ctx.JSON(http.StatusCreated, user_response.RegisterResponse{
		Message: "Success create user",
	})

}

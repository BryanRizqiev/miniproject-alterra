package user_controller

import (
	user_request "miniproject-alterra/module/user/controller/request"
	user_response "miniproject-alterra/module/user/controller/response"
	user_entity "miniproject-alterra/module/user/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService user_entity.UserServiceInterface
}

func NewUserController(userService user_entity.UserServiceInterface) *UserController {

	return &UserController{
		userService: userService,
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

	return ctx.JSON(http.StatusCreated, user_response.RegisterResponse{
		Message: "Success create user",
	})

}

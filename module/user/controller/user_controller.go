package user_controller

import (
	"miniproject-alterra/app/lib"
	"miniproject-alterra/app/validator"
	global_response "miniproject-alterra/module/global/controller/response"
	global_entity "miniproject-alterra/module/global/entity"
	user_request "miniproject-alterra/module/user/controller/request"
	user_response "miniproject-alterra/module/user/controller/response"
	user_entity "miniproject-alterra/module/user/entity"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService    user_entity.UserServiceInterface
	emailService   global_entity.EmailServiceInterface
	storageService global_entity.StorageServiceInterface
}

func NewUserController(userService user_entity.UserServiceInterface, emailService global_entity.EmailServiceInterface, storageService global_entity.StorageServiceInterface) *UserController {

	return &UserController{
		userService:    userService,
		emailService:   emailService,
		storageService: storageService,
	}

}

func (this *UserController) Register(ctx echo.Context) error {

	req := new(user_request.RegisterRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, user_response.StandartResponse{
			Message: "Request not valid",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}
	if req.DOB != "" {
		dob, err := time.Parse("2006-01-02", req.DOB)
		err = validator.DateValidation(err, dob)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, user_response.StandartResponse{
				Message: "Date of birth not valid.",
			})
		}
	}

	userDTO := user_entity.UserDTO{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Address:  req.Address,
		DOB:      req.DOB,
	}

	if err := this.userService.Register(userDTO); err != nil {
		err, ok := err.(*mysql.MySQLError)
		if ok && err.Number == 1062 {
			return ctx.JSON(http.StatusBadRequest, user_response.StandartResponse{
				Message: "Email is already in use.",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, user_response.StandartResponse{
			Message: "Server error",
		})
	}

	return ctx.JSON(http.StatusCreated, user_response.StandartResponse{
		Message: "Success create user",
	})

}

func (this *UserController) Login(ctx echo.Context) error {

	req := new(user_request.LoginRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, user_response.LoginResponse{
			Message: "Request not valid",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	userDTO := user_entity.UserDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := this.userService.Login(userDTO)
	if err != nil {
		errMessage := err.Error()
		if errMessage == "record not found" || errMessage == "credentials not valid" {
			return ctx.JSON(http.StatusBadRequest, user_response.LoginResponse{
				Message: "Credentials not valid",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, user_response.LoginResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, user_response.LoginResponse{
		Message: "Login success.",
		Token:   token,
	})

}

func (this *UserController) RequestVerified(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)
	err := this.userService.RequestVerification(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when request verification."
		errResStatus := http.StatusInternalServerError

		if errMessage == "email not verified" {
			errResMessage = "Email must be verified first."
			errResStatus = http.StatusBadRequest
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Request verified success.",
	})

}

func (this *UserController) VerifyEmail(ctx echo.Context) error {

	userId := ctx.Param("user-id")

	err := this.userService.Verify(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when request verification."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user already verified" {
			errResMessage = "User already verified."
			errResStatus = http.StatusBadRequest
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Verify user success.",
	})

}

func (this *UserController) RequestVerifyEmail(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)

	err := this.userService.SendVerifyEmail(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when request verification."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user already verified" {
			errResMessage = "User already verified."
			errResStatus = http.StatusBadRequest
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Request verified email success.",
	})

}

// Admin

func (this *UserController) GetRequestingUser(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)

	users, err := this.userService.GetRequestingUser(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when request verification."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		return ctx.JSON(errResStatus, user_response.GetRequestingUserRes{
			Message: errResMessage,
		})

	}

	var usersPres []user_response.UserPresentataion
	for _, user := range users {
		userDOB, _ := time.Parse(time.RFC3339, user.DOB.String)
		userPres := user_response.UserPresentataion{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			DOB:       userDOB.Format(time.DateOnly),
			Address:   user.Address.String,
			Phone:     user.Phone.String,
			Photo:     user.Photo.String,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
		}
		usersPres = append(usersPres, userPres)
	}

	return ctx.JSON(http.StatusOK, user_response.GetRequestingUserRes{
		Message: "Success get requesting users",
		Data:    usersPres,
	})

}

func (this *UserController) ChangeVerification(ctx echo.Context) error {

	req := new(user_request.ApproveVerificationReq)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, user_response.StandartResponse{
			Message: "Request not valid",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	userId, _ := lib.ExtractToken(ctx)

	err := this.userService.UpdateUserRole(userId, req.UserId, req.Role)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when change verification."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "record not found" {
			errResMessage = "User not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success change verification",
	})

}

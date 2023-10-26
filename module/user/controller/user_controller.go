package user_controller

import (
	"miniproject-alterra/app/lib"
	"miniproject-alterra/app/validator"
	"miniproject-alterra/module/dto"
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

func NewUserController(
	userService user_entity.UserServiceInterface,
	emailService global_entity.EmailServiceInterface,
	storageService global_entity.StorageServiceInterface,
) *UserController {

	return &UserController{
		userService:    userService,
		emailService:   emailService,
		storageService: storageService,
	}

}

func (this *UserController) Register(ctx echo.Context) error {

	req := new(user_request.RegisterRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
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
			return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
				Message: "Date of birth not valid.",
			})
		}
	}

	user := dto.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Address:  lib.NewNullString(req.Address),
		DOB:      lib.NewNullString(req.DOB),
		Phone:    lib.NewNullString(req.Phone),
	}

	err := this.userService.Register(user)
	if err != nil {

		err, ok := err.(*mysql.MySQLError)
		if ok && err.Number == 1062 {
			return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
				Message: "Email is already in use.",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, global_response.StandartResponse{
			Message: "Error when registration.",
		})

	}

	return ctx.JSON(http.StatusCreated, global_response.StandartResponse{
		Message: "Success registration.",
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

	user := dto.User{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := this.userService.Login(user)
	if err != nil {

		errMessage := err.Error()

		if errMessage == "record not found" || errMessage == "credentials not valid" {
			return ctx.JSON(http.StatusBadRequest, user_response.LoginResponse{
				Message: "Credentials not valid",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, user_response.LoginResponse{
			Message: "Error when login.",
		})

	}

	return ctx.JSON(http.StatusOK, user_response.LoginResponse{
		Message: "Login success.",
		Token:   token,
	})

}

func (this *UserController) RequestVerified(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)
	err := this.userService.RequestVerified(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when request verified."
		errResStatus := http.StatusInternalServerError

		if errMessage == "email not verified" {
			errResMessage = "Email must be verified first."
			errResStatus = http.StatusBadRequest
		}

		if errMessage == "record not found" {
			errResMessage = "Forbidden."
			errResStatus = http.StatusForbidden
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

	err := this.userService.VerifyEmail(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when verify email."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user already verified" {
			errResMessage = "User already verified."
			errResStatus = http.StatusBadRequest
		}

		if errMessage == "record not found" {
			errResMessage = "Forbidden."
			errResStatus = http.StatusForbidden
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Verify email success.",
	})

}

func (this *UserController) RequestVerifyEmail(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)

	err := this.userService.RequestVerifyEmail(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when request verify email."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user already verified" {
			errResMessage = "User already verified."
			errResStatus = http.StatusBadRequest
		}

		if errMessage == "record not found" {
			errResMessage = "Forbidden."
			errResStatus = http.StatusForbidden
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Request verify email success.",
	})

}

// Admin

func (this *UserController) GetRequestingUser(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)

	users, err := this.userService.GetRequestingUser(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when get requesting user."
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
			Id:        user.Id,
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

func (this *UserController) ChangeUserRole(ctx echo.Context) error {

	req := new(user_request.ApproveVerificationReq)
	userId := ctx.Param("user-id")

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "Request not valid",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	reqUserId, _ := lib.ExtractToken(ctx)

	err := this.userService.ChangeUserRole(reqUserId, userId, req.Role)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when change user role."
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
		Message: "Success change user role",
	})

}

func (this *UserController) GetAllUser(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)
	users, err := this.userService.GetAllUser(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when get users."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "record not found" {
			errResMessage = "Users not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponseWithData{
			Message: errResMessage,
		})

	}

	var userPresentations []user_response.UserPresentataion
	for _, user := range users {
		userDOB, _ := time.Parse(time.RFC3339, user.DOB.String)
		verifiedEmailAt := ""
		if user.VerifiedEmailAt.Valid {
			verifiedEmailAt = user.VerifiedEmailAt.Time.Format(lib.DATE_WITH_DAY_FORMAT)
		}

		userPresentataion := user_response.UserPresentataion{
			Id:              user.Id,
			Name:            user.Name,
			Email:           user.Email,
			DOB:             userDOB.Format(time.DateOnly),
			Address:         user.Address.String,
			Phone:           user.Phone.String,
			Photo:           user.Photo.String,
			Role:            user.Role,
			CreatedAt:       user.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
			RequestVerified: user.RequestVerified,
			VerifiedEmailAt: verifiedEmailAt,
		}
		userPresentations = append(userPresentations, userPresentataion)
	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponseWithData{
		Message: "Success get users.",
		Data:    userPresentations,
	})

}

func (this *UserController) UpdateUser(ctx echo.Context) error {

	req := new(user_request.UpdateRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
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
			return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
				Message: "Date of birth not valid.",
			})
		}
	}

	userId, _ := lib.ExtractToken(ctx)

	user := dto.User{
		Name:    req.Name,
		Address: lib.NewNullString(req.Address),
		DOB:     lib.NewNullString(req.DOB),
		Phone:   lib.NewNullString(req.Phone),
	}

	err := this.userService.UpdateUser(userId, user)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when update user."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Forbidden."
			errResStatus = http.StatusForbidden
		}

		return ctx.JSON(errResStatus, global_response.StandartResponseWithData{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success update user.",
	})

}

func (this *UserController) DeleteUser(ctx echo.Context) error {

	reqUserId, _ := lib.ExtractToken(ctx)
	userId := ctx.Param("user-id")

	err := this.userService.DeleteUser(reqUserId, userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when delete user."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "User not found."
			errResStatus = http.StatusNotFound
		}

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "nothing deleted" {
			errResMessage = "Nothing deleted."
			errResStatus = http.StatusBadRequest
		}

		return ctx.JSON(errResStatus, global_response.StandartResponseWithData{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success delete user.",
	})

}

func (this *UserController) UserSelfDelete(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)

	err := this.userService.UserSelfDelete(userId)
	if err != nil {

		errMessage := err.Error()
		errResMessage := "Error when user self delete."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Forbidden."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "nothing deleted" {
			errResMessage = "Nothing deleted."
			errResStatus = http.StatusBadRequest
		}

		return ctx.JSON(errResStatus, global_response.StandartResponseWithData{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success user self delete.",
	})

}

package user_controller

import (
	"miniproject-alterra/app/lib"
	"miniproject-alterra/app/validator"
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

	// format := global_entity.SendEmailFormat{
	// 	To:      userDTO.Email,
	// 	Cc:      userDTO.Email,
	// 	Subject: "Test",
	// 	Body:    "Registration Success",
	// }
	// // go this.emailService.SendEmail(format)

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

func (this *UserController) Verify(ctx echo.Context) error {

	userID, email := lib.ExtractToken(ctx)
	return ctx.JSON(http.StatusOK, user_response.VerifyResponse{
		Message: "Success",
		UserID:  userID,
		Email:   email,
	})

}

func (this *UserController) UploadPhoto(ctx echo.Context) error {

	file, err := ctx.FormFile("photo")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, user_response.StandartResponse{
			Message: "Request not valid.",
		})
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, user_response.StandartResponse{
			Message: "Error when read file.",
		})
	}
	defer src.Close()

	err = this.storageService.UploadFile("user-photo", file.Filename, src)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, user_response.StandartResponse{
			Message: "Error when upload in cloud storage.",
		})
	}

	return ctx.JSON(http.StatusOK, user_response.StandartResponse{
		Message: "Upload photo success.",
	})

}

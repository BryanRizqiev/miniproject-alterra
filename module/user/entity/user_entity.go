package user_entity

import (
	user_model "miniproject-alterra/module/user/repository/model"
	"time"
)

type UserDTO struct {
	ID              string
	Name            string
	Email           string
	Password        string
	DOB             string
	Address         string
	Phone           string
	VerifiedEmailAt time.Time
	Role            string
	RequestVerified string
	CreatedAt       time.Time
}

type (
	UserServiceInterface interface {
		Register(req UserDTO) error
		Login(req UserDTO) (string, error)
		GetAllUser() ([]UserDTO, error)
		RequestVerification(userId string) error
		Verify(userId string) error
		SendVerifyEmail(userId string) error
	}

	UserRepositoryInterface interface {
		InsertUser(userDTO UserDTO) error
		GetUserByEmail(email string) (UserDTO, error)
		GetAllUser() ([]UserDTO, error)
		UpdateUserRequestVerified(userId string) error
		UpdateUserVerifiedEmail(userId string) error
		CheckUserVerifiedEmail(userId string) (bool, error)
		FindUser(userId string) (user_model.User, error)
	}
)

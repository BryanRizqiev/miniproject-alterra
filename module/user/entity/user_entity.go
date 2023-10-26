package user_entity

import (
	"mime/multipart"
	"miniproject-alterra/module/dto"
)

type (
	UserServiceInterface interface {
		Register(user dto.User) error
		Login(user dto.User) (string, error)
		VerifyEmail(userId string) error

		RequestVerified(userId string) error
		RequestVerifyEmail(userId string) error

		GetAllUser(userId string) ([]dto.User, error)
		GetRequestingUser(userId string) ([]dto.User, error)
		ChangeUserRole(reqUserId string, userId string, role string) error
		UpdateUser(userId string, payload dto.User) error
		DeleteUser(reqUserId, userId string) error
		UserSelfDelete(userId string) error
		UpdatePhoto(userId, filename string, image multipart.File) error
		GetUserProfile(userId string) (dto.User, error)
	}

	UserRepositoryInterface interface {
		InsertUser(user dto.User) error
		GetUserByEmail(email string) (dto.User, error)
		GetAllUser() ([]dto.User, error)

		UpdateUserRequestVerified(userId string) error
		UpdateUserVerifiedEmail(userId string) error
		CheckUserVerifiedEmail(userId string) (bool, error)

		UpdateUserRole(userId string, role string) error
		FindUser(userId string) (dto.User, error)
		GetRequestingUser() ([]dto.User, error)
		UpdateUser(user dto.User) error
		DeleteUser(user dto.User) error
		UpdatePhoto(fileName string, user dto.User) error
	}
)

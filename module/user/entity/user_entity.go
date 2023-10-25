package user_entity

import (
	"database/sql"
	"miniproject-alterra/module/dto"
	"time"

	"gorm.io/gorm"
)

type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	*Base

	ID              string `gorm:"primaryKey"`
	Name            string
	Email           string
	Password        string
	DOB             sql.NullString `gorm:"column:dob"`
	Address         sql.NullString
	Phone           sql.NullString
	Photo           sql.NullString
	VerifiedEmailAt sql.NullTime
	Role            string `gorm:"default:user"`
	RequestVerified string `gorm:"default:default"`
}

type (
	UserServiceInterface interface {
		Register(user dto.User) error
		Login(user dto.User) (string, error)

		GetAllUser() ([]dto.User, error)
		RequestVerified(userId string) error
		VerifyEmail(userId string) error

		RequestVerifyEmail(userId string) error
		GetRequestingUser(userId string) ([]dto.User, error)
		ChangeUserRole(reqUserId string, userId string, role string) error
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
	}
)

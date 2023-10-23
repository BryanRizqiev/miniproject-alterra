package user_entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
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
		Register(req UserDTO) error
		Login(req UserDTO) (string, error)
		GetAllUser() ([]UserDTO, error)
		RequestVerification(userId string) error
		Verify(userId string) error
		SendVerifyEmail(userId string) error
		GetRequestingUser(userId string) ([]User, error)
		UpdateUserRole(reqUserId string, userId string, role string) error
	}

	UserRepositoryInterface interface {
		InsertUser(userDTO UserDTO) error
		GetUserByEmail(email string) (UserDTO, error)
		GetAllUser() ([]UserDTO, error)
		UpdateUserRequestVerified(userId string) error
		UpdateUserVerifiedEmail(userId string) error
		CheckUserVerifiedEmail(userId string) (bool, error)
		FindUser(userId string) (User, error)
		GetRequestingUser() ([]User, error)
		UpdateUserRole(userId string, role string) error
	}
)

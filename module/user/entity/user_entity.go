package user_entity

import "time"

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
		Login(req UserDTO) (UserDTO error)
		GetAllUser() ([]UserDTO, error)
	}

	UserRepositoryInterface interface {
		InsertUser(userDTO UserDTO) error
		CheckUser(userDTO UserDTO) (UserDTO error)
		GetAllUser() ([]UserDTO, error)
	}
)

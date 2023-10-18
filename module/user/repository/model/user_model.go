package user_model

import (
	"database/sql"
	"miniproject-alterra/app/lib"
)

type User struct {
	*lib.Base

	ID              string
	Name            string
	Email           string
	Password        string
	Address         sql.NullString
	DOB             sql.NullString `gorm:"column:dob"`
	Phone           sql.NullString
	VerifiedEmailAt sql.NullTime
	Role            sql.NullString
	RequestVerified sql.NullString
	CreatedAt       sql.NullTime
}

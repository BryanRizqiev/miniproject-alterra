package user_model

import (
	"database/sql"
	"miniproject-alterra/app/lib"
)

type User struct {
	*lib.Base

	Id              string `gorm:"primaryKey"`
	Name            string
	Email           string
	Password        string
	Address         sql.NullString
	DOB             sql.NullString `gorm:"column:dob"`
	Phone           sql.NullString
	VerifiedEmailAt sql.NullTime
	Role            string `gorm:"default:user"`
	RequestVerified string `gorm:"default:default"`
}

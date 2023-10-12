package user_model

import (
	"database/sql"
	"miniproject-alterra/app/lib"
)

type User struct {
	*lib.Base

	Email    string `gorm:"uniqueIndex"`
	Name     string
	Password string
	Address  sql.NullString
}

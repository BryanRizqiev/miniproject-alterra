package event_model

import (
	"database/sql"
	"miniproject-alterra/app/lib"
	user_model "miniproject-alterra/module/user/repository/model"
)

type Event struct {
	*lib.Base

	ID                string
	Title             string
	Location          string
	LocationURL       sql.NullString `gorm:"column:location_url"`
	Description       sql.NullString
	Image             sql.NullString
	Status            string `gorm:"default:waiting"`
	RecommendedAction sql.NullString
	UserID            string
	CreatedBy         user_model.User `gorm:"foreignKey:UserID"`
}

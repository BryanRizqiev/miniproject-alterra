package dto

import (
	"database/sql"
	"miniproject-alterra/app/lib"
)

type Evidence struct {
	*lib.Base

	Id        string
	Content   string
	Image     string
	CreatedBy string
	User      User `gorm:"foreignKey:CreatedBy"`
	EventId   string
	Event     Event `gorm:"foreignKey:EventId"`
}

type User struct {
	*lib.Base

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

type Event struct {
	*lib.Base

	Id                string `gorm:"primaryKey"`
	Title             string
	Location          string
	LocationURL       sql.NullString `gorm:"column:location_url"`
	Description       sql.NullString
	Image             sql.NullString
	Status            string `gorm:"default:waiting"`
	RecommendedAction sql.NullString
	UserId            string
	CreatedBy         User       `gorm:"foreignKey:UserId"`
	Evidences         []Evidence `gorm:"foreignKey:EventId"`
}

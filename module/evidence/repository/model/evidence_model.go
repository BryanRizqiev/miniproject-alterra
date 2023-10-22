package evd_model

import (
	"miniproject-alterra/app/lib"
	event_model "miniproject-alterra/module/events/repository/model"
	user_model "miniproject-alterra/module/user/repository/model"
)

type Evidence struct {
	*lib.Base

	ID        string
	Content   string
	Image     string
	CreatedBy string
	User      user_model.User `gorm:"foreignKey:CreatedBy"`
	EventId   string
	Event     event_model.Event `gorm:"foreignKey:EventId"`
}

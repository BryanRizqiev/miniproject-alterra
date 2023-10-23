package event_entity

import (
	"mime/multipart"
	"miniproject-alterra/module/dto"
	user_model "miniproject-alterra/module/user/repository/model"
	"time"
)

type EventDTO struct {
	ID                string
	Title             string
	Location          string
	LocationURL       string
	Description       string
	Image             string
	Status            string
	RecommendedAction string
	UserID            string
	CreatedBy         user_model.User
	CreatedAt         time.Time
}

type (
	IEventReposistory interface {
		InsertEvent(eventDTO EventDTO) error
		GetEvent() ([]dto.Event, error)
	}

	IEventService interface {
		CreateEvent(userID string, evtDTO EventDTO, image multipart.File) error
		GetEvent() ([]dto.Event, error)
	}
)

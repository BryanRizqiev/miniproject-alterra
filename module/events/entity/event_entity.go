package event_entity

import (
	"mime/multipart"
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
		GetEvent() ([]EventDTO, error)
	}

	IEventService interface {
		CreateEvent(userID string, evtDTO EventDTO, image multipart.File) error
		GetEvent() ([]EventDTO, error)
	}
)

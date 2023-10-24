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
		UpdateRecommendedAction(evt dto.Event, value string) error
		InsertEvent(evtD EventDTO) (dto.Event, error)
		GetEvent() ([]dto.Event, error)
		UpdateEventStatus(event dto.Event, status string) error
		FindEvent(eventId string) (dto.Event, error)
		FindOwnEvent(userId, eventId string) (dto.Event, error)
		UpdateEvent(event dto.Event) (dto.Event, error)
	}

	IEventService interface {
		CreateEvent(userID string, evtDTO EventDTO, image multipart.File) error
		GetEvent() ([]dto.Event, error)
		PublishEvent(userId, evtId string) error
		UpdateEvent(userId, eventId string, payload dto.Event) error
	}
)

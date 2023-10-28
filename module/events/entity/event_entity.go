package event_entity

import (
	"mime/multipart"
	"miniproject-alterra/module/dto"
	user_model "miniproject-alterra/module/user/repository/model"
	"time"
)

type EventDTO struct {
	Id                string
	Title             string
	Location          string
	LocationURL       string
	Description       string
	Image             string
	Status            string
	RecommendedAction string
	UserId            string
	CreatedBy         user_model.User
	CreatedAt         time.Time
}

type (
	IEventReposistory interface {
		InsertEvent(eventD EventDTO) (dto.Event, error)
		UpdateRecommendedAction(event dto.Event, value string) error
		GetEvent() ([]dto.Event, error)
		UpdateEventStatus(event dto.Event, status string) error
		FindEvent(eventId string) (dto.Event, error)
		FindOwnEvent(userId, eventId string) (dto.Event, error)
		UpdateEvent(event dto.Event) (dto.Event, error)
		GetWaitingEvents() ([]dto.Event, error)
		DeleteEvent(event dto.Event) error
		UpdateImage(fileName string, event dto.Event) error
		GetAllEvent() ([]dto.Event, error)
	}

	IEventService interface {
		CreateEvent(userId string, eventDTO EventDTO, image multipart.File) error
		GetEvent() ([]dto.Event, error)
		PublishEvent(userId, eventId string) error
		UpdateEvent(userId, eventId string, payload dto.Event) error
		GetWaitingEvents(userId string) ([]dto.Event, error)
		UpdateImage(userId, eventId, filename string, image multipart.File) error
		DeleteEvent(userId, eventId string) error
		TakedownEvent(userId, eventId string) error
		GetAllEvent(userId string) ([]dto.Event, error)
	}
)

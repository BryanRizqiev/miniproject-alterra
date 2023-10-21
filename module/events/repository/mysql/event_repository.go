package mysql_event_repository

import (
	"miniproject-alterra/app/lib"
	event_entity "miniproject-alterra/module/events/entity"
	event_model "miniproject-alterra/module/events/repository/model"

	"gorm.io/gorm"
)

type EventReposistory struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) event_entity.IEventReposistory {

	return &EventReposistory{
		db: db,
	}

}

func (this *EventReposistory) InsertEvent(eventDTO event_entity.EventDTO) error {

	event := event_model.Event{
		ID:          lib.NewUuid(),
		Title:       eventDTO.Title,
		Location:    eventDTO.Location,
		LocationURL: lib.NewNullString(eventDTO.LocationURL),
		Description: lib.NewNullString(eventDTO.Location),
		UserID:      eventDTO.UserID,
	}

	tx := this.db.Create(&event)

	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

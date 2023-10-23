package mysql_event_repository

import (
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
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

func (this *EventReposistory) InsertEvent(evtD event_entity.EventDTO) error {

	evtM := event_model.Event{
		ID:          lib.NewUuid(),
		Title:       evtD.Title,
		Location:    evtD.Location,
		LocationURL: lib.NewNullString(evtD.LocationURL),
		Description: lib.NewNullString(evtD.Location),
		UserID:      evtD.UserID,
		Image:       lib.NewNullString(evtD.Image),
	}

	tx := this.db.Create(&evtM)

	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (this *EventReposistory) GetEvent() ([]dto.Event, error) {

	var evts []dto.Event

	err := this.db.Where("status = ?", "publish").Preload("CreatedBy").Preload("Evidences.User").Find(&evts).Error
	if err != nil {
		return []dto.Event{}, err
	}

	return evts, nil

}

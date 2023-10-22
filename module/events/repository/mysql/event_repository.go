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

func (this *EventReposistory) GetEvent() ([]event_entity.EventDTO, error) {

	var evtsM []event_model.Event
	var evtsD []event_entity.EventDTO

	tx := this.db.Where("status = ?", "publish").Preload("CreatedBy").Find(&evtsM)
	if tx.Error != nil {
		return evtsD, tx.Error
	}

	for _, value := range evtsM {
		evtDTO := event_entity.EventDTO{
			ID:                value.ID,
			Title:             value.Title,
			Location:          value.Location,
			LocationURL:       value.LocationURL.String,
			Description:       value.Description.String,
			Image:             value.Image.String,
			RecommendedAction: value.RecommendedAction.String,
			CreatedBy:         value.CreatedBy,
		}
		evtsD = append(evtsD, evtDTO)
	}

	return evtsD, nil

}

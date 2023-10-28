package mysql_event_repository

import (
	"errors"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
	event_entity "miniproject-alterra/module/events/entity"

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

func (this *EventReposistory) InsertEvent(eventDTO event_entity.EventDTO) (dto.Event, error) {

	event := dto.Event{
		Id:          lib.NewUuid(),
		Title:       eventDTO.Title,
		Location:    eventDTO.Location,
		LocationURL: lib.NewNullString(eventDTO.LocationURL),
		Description: lib.NewNullString(eventDTO.Description),
		UserId:      eventDTO.UserId,
		Image:       lib.NewNullString(eventDTO.Image),
		Status:      lib.InsertDefaultValue("waiting", eventDTO.Status),
	}

	tx := this.db.Create(&event)

	if tx.Error != nil {
		return dto.Event{}, tx.Error
	}

	return event, nil

}

func (this *EventReposistory) GetEvent() ([]dto.Event, error) {

	var evts []dto.Event

	err := this.db.Where("status = ?", "publish").Preload("CreatedBy").Preload("Evidences.User").Find(&evts).Error
	if err != nil {
		return []dto.Event{}, err
	}

	return evts, nil

}

func (this *EventReposistory) GetAllEvent() ([]dto.Event, error) {

	var event []dto.Event

	err := this.db.Unscoped().Preload("CreatedBy").Find(&event).Error
	if err != nil {
		return []dto.Event{}, err
	}

	return event, nil

}

func (this *EventReposistory) UpdateRecommendedAction(evt dto.Event, value string) error {

	err := this.db.Model(&evt).Update("recommended_action", value).Error
	if err != nil {
		return err
	}

	return nil

}

func (this *EventReposistory) FindEvent(eventId string) (dto.Event, error) {

	var event dto.Event

	tx := this.db.First(&event, "id = ?", eventId)
	if tx.Error != nil {
		return dto.Event{}, tx.Error
	}

	return event, nil

}

func (this *EventReposistory) UpdateEventStatus(event dto.Event, status string) error {

	err := this.db.First(&event, "id = ?", event.Id).Error
	if err != nil {
		return err
	}

	err = this.db.Model(&event).Update("status", status).Error
	if err != nil {
		return err
	}

	return nil

}

func (this *EventReposistory) GetWaitingEvents() ([]dto.Event, error) {

	var events []dto.Event

	err := this.db.Where("status = ?", "waiting").Preload("CreatedBy").Find(&events).Error
	if err != nil {
		return []dto.Event{}, err
	}

	return events, nil

}

func (this *EventReposistory) FindOwnEvent(userId, eventId string) (dto.Event, error) {

	var event dto.Event

	tx := this.db.Where("user_id", userId).First(&event, "id = ?", eventId)
	if tx.Error != nil {
		return dto.Event{}, tx.Error
	}

	return event, nil

}

func (this *EventReposistory) UpdateEvent(event dto.Event) (dto.Event, error) {

	err := this.db.Save(&event).Error
	if err != nil {
		return dto.Event{}, err
	}

	return event, nil

}

func (this *EventReposistory) UpdateImage(fileName string, event dto.Event) error {

	err := this.db.Model(&event).Update("image", fileName).Error
	if err != nil {
		return err
	}

	return nil

}

func (this *EventReposistory) DeleteEvent(event dto.Event) error {

	tx := this.db.Delete(&event)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected < 1 {
		return errors.New("nothing deleted")
	}

	return nil

}

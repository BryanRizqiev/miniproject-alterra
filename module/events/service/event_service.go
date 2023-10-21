package event_service

import (
	event_entity "miniproject-alterra/module/events/entity"
)

type EventService struct {
	evtRepo event_entity.IEventReposistory
}

func NewEventService(evtRepo event_entity.IEventReposistory) event_entity.IEventService {

	return &EventService{
		evtRepo: evtRepo,
	}

}

func (this *EventService) CreateEvent(userID string, evtDTO event_entity.EventDTO) error {

	evtDTO.UserID = userID

	err := this.evtRepo.InsertEvent(evtDTO)
	if err != nil {
		return err
	}

	return nil

}

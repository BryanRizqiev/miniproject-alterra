package event_service

import (
	"fmt"
	"mime/multipart"
	"miniproject-alterra/app/lib"
	event_entity "miniproject-alterra/module/events/entity"
	global_entity "miniproject-alterra/module/global/entity"
	user_model "miniproject-alterra/module/user/repository/model"
	"path/filepath"
	"strings"
)

type EventService struct {
	evtRepo    event_entity.IEventReposistory
	storageSvc global_entity.StorageServiceInterface
}

func NewEventService(evtRepo event_entity.IEventReposistory, storageSvc global_entity.StorageServiceInterface) event_entity.IEventService {

	return &EventService{
		evtRepo:    evtRepo,
		storageSvc: storageSvc,
	}

}

func (this *EventService) CreateEvent(userID string, evtD event_entity.EventDTO, image multipart.File) error {

	evtD.UserID = userID

	fileExt := strings.ToLower(filepath.Ext(evtD.Image))
	newFilename := fmt.Sprintf("%s-%s%s", "event", lib.RandomString(8), fileExt)
	evtD.Image = newFilename

	var err error
	err = this.storageSvc.UploadFile("event", newFilename, image)
	if err != nil {
		return err
	}
	err = this.evtRepo.InsertEvent(evtD)
	if err != nil {
		return err
	}

	return nil

}

func (this *EventService) GetEvent() ([]event_entity.EventDTO, error) {

	evtsD, err := this.evtRepo.GetEvent()
	for i := range evtsD {
		usr := user_model.User{}
		usr.Name = evtsD[i].CreatedBy.Name
		usr.Role = evtsD[i].CreatedBy.Role
		evtsD[i].CreatedBy = usr

		url, errURL := this.storageSvc.GetUrl("event", evtsD[i].Image)
		if errURL != nil {
			return []event_entity.EventDTO{}, err
		}
		evtsD[i].Image = url
	}

	if err != nil {
		return []event_entity.EventDTO{}, err
	}

	return evtsD, nil

}

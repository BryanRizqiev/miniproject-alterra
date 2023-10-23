package event_service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
	event_entity "miniproject-alterra/module/events/entity"
	global_entity "miniproject-alterra/module/global/entity"
	global_service "miniproject-alterra/module/global/service"
	"path/filepath"
	"strings"
)

type EventService struct {
	evtRepo    event_entity.IEventReposistory
	storageSvc global_entity.StorageServiceInterface
	openai     *global_service.OpenAIService
	globalRepo global_entity.IGlobalRepository
}

func NewEventService(
	evtRepo event_entity.IEventReposistory,
	storageSvc global_entity.StorageServiceInterface,
	openai *global_service.OpenAIService,
	globalRepo global_entity.IGlobalRepository) event_entity.IEventService {

	return &EventService{
		evtRepo:    evtRepo,
		storageSvc: storageSvc,
		openai:     openai,
		globalRepo: globalRepo,
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

	user, err := this.globalRepo.GetUser(evtD.UserID)
	if err != nil {
		return err
	}
	if user.Role != "user" {
		evtD.Status = "publish"
	}

	evt, err := this.evtRepo.InsertEvent(evtD)

	go this.updateRecommendedAction(evt)

	if err != nil {
		return err
	}

	return nil

}

func (this *EventService) updateRecommendedAction(evt dto.Event) {

	promt := fmt.Sprintf("%s, %s", evt.Title, evt.Description.String)
	rec, err := this.openai.GetRecommendedAction(promt)
	if err != nil {
		panic(err.Error())
	}

	err = this.evtRepo.UpdateRecommendedAction(evt, rec)
	if err != nil {
		panic(err.Error())
	}

}

func (this *EventService) GetEvent() ([]dto.Event, error) {

	evts, err := this.evtRepo.GetEvent()
	for i := range evts {
		usr := dto.User{}
		usr.Name = evts[i].CreatedBy.Name
		usr.Role = evts[i].CreatedBy.Role
		evts[i].CreatedBy = usr

		url, errURL := this.storageSvc.GetUrl("event", evts[i].Image.String)
		if errURL != nil {
			return []dto.Event{}, err
		}
		evts[i].Image.String = url

		evds := evts[i].Evidences
		for j := range evds {
			url, errURL := this.storageSvc.GetUrl("event", evts[i].Image.String)
			if errURL != nil {
				return []dto.Event{}, err
			}
			evds[j].Image = url
		}

		evts[i].Evidences = evds
	}

	if err != nil {
		return []dto.Event{}, err
	}

	return evts, nil

}

func (this *EventService) PublishEvent(userId, evtId string) error {

	user, err := this.globalRepo.GetUser(userId)
	if err != nil {
		return err
	}

	if !lib.CheckIsAdmin(user) {
		return errors.New("user not allowed")
	}

	event, err := this.evtRepo.FindEvent(evtId)
	if err != nil {
		return err
	}
	err = this.evtRepo.UpdateEventStatus(event, "publish")
	if err != nil {
		return err
	}

	return nil

}

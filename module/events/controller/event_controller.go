package event_controller

import (
	"fmt"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/app/validator"
	"miniproject-alterra/module/dto"
	event_request "miniproject-alterra/module/events/controller/request"
	evt_response "miniproject-alterra/module/events/controller/response"
	event_entity "miniproject-alterra/module/events/entity"
	evd_res "miniproject-alterra/module/evidence/controller/response"
	global_response "miniproject-alterra/module/global/controller/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EventController struct {
	eventSvc event_entity.IEventService
}

func NewEventController(eventSvc event_entity.IEventService) *EventController {

	return &EventController{
		eventSvc: eventSvc,
	}

}

func (this *EventController) CreateEvent(ctx echo.Context) error {

	req := new(event_request.CreateEvtReq)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "Request not valid.",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "image required.",
		})
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, global_response.StandartResponse{
			Message: "Error when read file.",
		})
	}
	defer src.Close()

	if !validator.ImageValidation(file) {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "image must be valid image.",
		})
	}

	if req.LocationURL != "" && !validator.GoogleMapsURLValidator(req.LocationURL) {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "location_url not valid.",
		})
	}

	userID, _ := lib.ExtractToken(ctx)
	evtDTO := event_entity.EventDTO{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		LocationURL: req.LocationURL,
		Image:       file.Filename,
	}

	err = this.eventSvc.CreateEvent(userID, evtDTO, src)
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, global_response.StandartResponse{
			Message: "Server error.",
		})
	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success create event.",
	})

}

func (this *EventController) PublishEvent(ctx echo.Context) error {

	eventId := ctx.Param("event-id")
	userId, _ := lib.ExtractToken(ctx)

	err := this.eventSvc.PublishEvent(userId, eventId)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when publish event."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "record not found" {
			errResMessage = "Event not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success publish event.",
	})

}

func (this *EventController) TakedownEvent(ctx echo.Context) error {

	eventId := ctx.Param("event-id")
	userId, _ := lib.ExtractToken(ctx)

	err := this.eventSvc.TakedownEvent(userId, eventId)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when takedown event."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "record not found" {
			errResMessage = "Event not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success takedown event.",
	})

}

func (this *EventController) GetEvent(ctx echo.Context) error {

	evts, err := this.eventSvc.GetEvent()
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, evt_response.GetEventResponse{
			Message: "Error when get event.",
		})
	}

	var presentations []evt_response.EventPresentation
	for _, value := range evts {
		verfied := true
		if value.CreatedBy.Role == "user" {
			verfied = false
		}

		var evdPresentations []evd_res.EvdsPresentation
		for _, evd := range value.Evidences {
			evdVerified := true
			if evd.User.Role == "user" {
				evdVerified = false
			}
			evdPresentation := evd_res.EvdsPresentation{
				Content:   evd.Content,
				Image:     evd.Image,
				CreatedAt: evd.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
				CreatedBy: evd.User.Name,
				Verified:  evdVerified,
			}
			evdPresentations = append(evdPresentations, evdPresentation)
		}

		presentation := evt_response.EventPresentation{
			Id:                value.Id,
			Title:             value.Title,
			Location:          value.Location,
			LocationURL:       value.LocationURL.String,
			Description:       value.Description.String,
			Image:             value.Image.String,
			RecommendedAction: value.RecommendedAction.String,
			CreatedBy:         value.CreatedBy.Name,
			Verified:          verfied,
			Evidences:         evdPresentations,
			CreatedAt:         value.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
		}
		presentations = append(presentations, presentation)
	}

	return ctx.JSON(http.StatusOK, evt_response.GetEventResponse{
		Message: "Success get event.",
		Data:    presentations,
	})

}

func (this *EventController) GetAllEvent(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)
	events, err := this.eventSvc.GetAllEvent(userId)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when get all event."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		return ctx.JSON(errResStatus, evt_response.GetEventResponse{
			Message: errResMessage,
		})

	}

	var eventPresentations []evt_response.EventPresentation
	for _, event := range events {
		verfied := true
		if event.CreatedBy.Role == "user" {
			verfied = false
		}

		eventPresentation := evt_response.EventPresentation{
			Id:                event.Id,
			Title:             event.Title,
			Location:          event.Location,
			LocationURL:       event.LocationURL.String,
			Description:       event.Description.String,
			Image:             event.Image.String,
			RecommendedAction: event.RecommendedAction.String,
			CreatedBy:         event.CreatedBy.Name,
			Verified:          verfied,
			CreatedAt:         event.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
		}

		eventPresentations = append(eventPresentations, eventPresentation)
	}

	return ctx.JSON(http.StatusOK, evt_response.GetEventResponse{
		Message: "Success get waiting all event.",
		Data:    eventPresentations,
	})

}

func (this *EventController) GetWaitingEvents(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)
	events, err := this.eventSvc.GetWaitingEvents(userId)

	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when get waiting events."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		return ctx.JSON(errResStatus, evt_response.GetEventResponse{
			Message: errResMessage,
		})

	}

	var eventPresentations []evt_response.EventPresentation
	for _, event := range events {
		verfied := true
		if event.CreatedBy.Role == "user" {
			verfied = false
		}

		eventPresentation := evt_response.EventPresentation{
			Id:                event.Id,
			Title:             event.Title,
			Location:          event.Location,
			LocationURL:       event.LocationURL.String,
			Description:       event.Description.String,
			Image:             event.Image.String,
			RecommendedAction: event.RecommendedAction.String,
			CreatedBy:         event.CreatedBy.Name,
			Verified:          verfied,
			CreatedAt:         event.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
		}

		eventPresentations = append(eventPresentations, eventPresentation)
	}

	return ctx.JSON(http.StatusOK, evt_response.GetEventResponse{
		Message: "Success get waiting events.",
		Data:    eventPresentations,
	})

}

func (this *EventController) UpdateEvent(ctx echo.Context) error {

	eventId := ctx.Param("event-id")
	req := new(event_request.UpdateEventReq)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "Request not valid.",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	if req.LocationURL != "" && !validator.GoogleMapsURLValidator(req.LocationURL) {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "location_url not valid.",
		})
	}

	userId, _ := lib.ExtractToken(ctx)
	event := dto.Event{
		Title:       req.Title,
		Location:    req.Location,
		LocationURL: lib.NewNullString(req.LocationURL),
		Description: lib.NewNullString(req.Description),
	}

	err := this.eventSvc.UpdateEvent(userId, eventId, event)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when update event."
		errResStatus := http.StatusInternalServerError

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		if errMessage == "record not found" {
			errResMessage = "Event not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success update event.",
	})

}

func (this *EventController) UpdateImage(ctx echo.Context) error {

	userId, _ := lib.ExtractToken(ctx)
	eventId := ctx.Param("event-id")

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "image required.",
		})
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, global_response.StandartResponse{
			Message: "Error when read file.",
		})
	}

	err = this.eventSvc.UpdateImage(userId, eventId, file.Filename, src)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when update event image."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Event not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success update event image.",
	})

}

func (this *EventController) DeleteEvent(ctx echo.Context) error {

	eventId := ctx.Param("event-id")
	userId, _ := lib.ExtractToken(ctx)

	err := this.eventSvc.DeleteEvent(userId, eventId)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when delete event."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Event not found."
			errResStatus = http.StatusNotFound
		}

		if errMessage == "nothing deleted" {
			errResMessage = "Nothing deleted."
			errResStatus = http.StatusBadRequest
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success delete event.",
	})

}

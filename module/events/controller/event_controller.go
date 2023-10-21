package event_controller

import (
	"miniproject-alterra/app/lib"
	event_request "miniproject-alterra/module/events/controller/request"
	event_entity "miniproject-alterra/module/events/entity"
	global_response "miniproject-alterra/module/global/controller/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EventController struct {
	evtSvc event_entity.IEventService
}

func NewEventController(evtSvc event_entity.IEventService) *EventController {

	return &EventController{
		evtSvc: evtSvc,
	}

}

func (this *EventController) CreateEvent(ctx echo.Context) error {

	req := new(event_request.CreateEvtReq)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "Request not valid",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	// reg, err := regexp.Compile()

	userID, _ := lib.ExtractToken(ctx)
	evtDTO := event_entity.EventDTO{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		LocationURL: req.LocationURL,
	}

	err := this.evtSvc.CreateEvent(userID, evtDTO)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, global_response.StandartResponse{
			Message: "Server error.",
		})
	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success create event.",
	})

}

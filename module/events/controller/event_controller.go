package event_controller

import (
	"miniproject-alterra/app/lib"
	"miniproject-alterra/app/validator"
	event_request "miniproject-alterra/module/events/controller/request"
	evt_response "miniproject-alterra/module/events/controller/response"
	event_entity "miniproject-alterra/module/events/entity"
	evd_res "miniproject-alterra/module/evidence/controller/response"
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
			Message: "Request not valid.",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "Request not valid.",
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

	err = this.evtSvc.CreateEvent(userID, evtDTO, src)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, global_response.StandartResponse{
			Message: "Server error.",
		})
	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success create event.",
	})

}

func (this *EventController) GetEvent(ctx echo.Context) error {

	evts, err := this.evtSvc.GetEvent()
	if err != nil {
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
		}
		presentations = append(presentations, presentation)
	}

	return ctx.JSON(http.StatusOK, evt_response.GetEventResponse{
		Message: "Success get event.",
		Data:    presentations,
	})

}

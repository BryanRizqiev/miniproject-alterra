package evd_controller

import (
	"miniproject-alterra/app/lib"
	"miniproject-alterra/app/validator"
	evd_req "miniproject-alterra/module/evidence/controller/request"
	evd_res "miniproject-alterra/module/evidence/controller/response"
	evd_entity "miniproject-alterra/module/evidence/entity"
	global_response "miniproject-alterra/module/global/controller/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EvidenceController struct {
	evdSvc evd_entity.IEvidenceService
}

func NewEvidenceController(evdSvc evd_entity.IEvidenceService) *EvidenceController {

	return &EvidenceController{
		evdSvc: evdSvc,
	}

}

func (this *EvidenceController) CreateEvidence(ctx echo.Context) error {

	req := new(evd_req.CreateEvdReq)

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

	userId, _ := lib.ExtractToken(ctx)
	evdD := evd_entity.EvidenceDTO{
		Content: req.Content,
		Image:   file.Filename,
	}

	err = this.evdSvc.CreateEvidence(userId, req.EventId, src, evdD)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, global_response.StandartResponse{
			Message: "Error when create evidence.",
		})
	}

	return ctx.JSON(http.StatusCreated, global_response.StandartResponse{
		Message: "Success create evidence.",
	})

}

func (this *EvidenceController) GetEvidences(ctx echo.Context) error {

	eventId := ctx.Param("event-id")

	evdsD, err := this.evdSvc.GetEvidences(eventId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, evd_res.GetEvdsRes{
			Message: "Error when get evidences.",
		})
	}

	var evdsPresentator []evd_res.EvdsPresentator
	for _, value := range evdsD {
		evdPresentator := evd_res.EvdsPresentator{
			Content:   value.Content,
			Image:     value.Image,
			CreatedAt: value.CreatedAt,
		}
		evdsPresentator = append(evdsPresentator, evdPresentator)
	}

	return ctx.JSON(http.StatusInternalServerError, evd_res.GetEvdsRes{
		Message: "Success get evidences.",
		Data:    evdsPresentator,
	})

}

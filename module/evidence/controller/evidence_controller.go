package evd_controller

import (
	"miniproject-alterra/app/lib"
	"miniproject-alterra/app/validator"
	"miniproject-alterra/module/dto"
	evd_req "miniproject-alterra/module/evidence/controller/request"
	evd_res "miniproject-alterra/module/evidence/controller/response"
	evd_entity "miniproject-alterra/module/evidence/entity"
	global_response "miniproject-alterra/module/global/controller/response"
	"net/http"

	"github.com/go-sql-driver/mysql"
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

	userId, _ := lib.ExtractToken(ctx)
	evidence := dto.Evidence{
		Content: req.Content,
		Image:   file.Filename,
	}

	err = this.evdSvc.CreateEvidence(userId, req.EventId, src, evidence)
	if err != nil {

		err, ok := err.(*mysql.MySQLError)
		if ok && err.Number == 1452 {
			return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
				Message: "Event not found.",
			})
		}

		errResMessage := "Error when get create evidence."
		errResStatus := http.StatusInternalServerError

		return ctx.JSON(errResStatus, global_response.StandartResponseWithData{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusCreated, global_response.StandartResponse{
		Message: "Success create evidence.",
	})

}

func (this *EvidenceController) GetEvidences(ctx echo.Context) error {

	eventId := ctx.Param("event-id")

	evidences, err := this.evdSvc.GetEvidences(eventId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, evd_res.GetEvdsRes{
			Message: "Error when get evidences.",
		})
	}

	var evdsPresentator []evd_res.EvdsPresentation
	for _, evidence := range evidences {
		evdPresentator := evd_res.EvdsPresentation{
			Content:   evidence.Content,
			Image:     evidence.Image,
			CreatedAt: evidence.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
		}
		evdsPresentator = append(evdsPresentator, evdPresentator)
	}

	return ctx.JSON(http.StatusInternalServerError, evd_res.GetEvdsRes{
		Message: "Success get evidences.",
		Data:    evdsPresentator,
	})

}

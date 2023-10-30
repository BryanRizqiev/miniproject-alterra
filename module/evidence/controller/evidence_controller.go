package evd_controller

import (
	"fmt"
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

		fmt.Println(err.Error())
		err, ok := err.(*mysql.MySQLError)
		if ok && err.Number == 1452 {
			return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
				Message: "Event not found.",
			})
		}

		errResMessage := "Error when create evidence."
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
	userId, _ := lib.ExtractToken(ctx)

	evidences, err := this.evdSvc.GetEvidences(userId, eventId)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when get evidences."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Evidences not found."
			errResStatus = http.StatusNotFound
		}

		if errMessage == "user not allowed" {
			errResMessage = "User not allowed."
			errResStatus = http.StatusForbidden
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	var evdsPresentator []evd_res.EvdsPresentation
	for _, evidence := range evidences {
		isVerified := false
		if evidence.User.Role != "user" {
			isVerified = true
		}
		evdPresentator := evd_res.EvdsPresentation{
			Content:   evidence.Content,
			Image:     evidence.Image,
			CreatedAt: evidence.CreatedAt.Format(lib.DATE_WITH_DAY_FORMAT),
			CreatedBy: evidence.User.Name,
			Verified:  isVerified,
		}
		evdsPresentator = append(evdsPresentator, evdPresentator)
	}

	return ctx.JSON(http.StatusInternalServerError, evd_res.GetEvdsRes{
		Message: "Success get evidences.",
		Data:    evdsPresentator,
	})

}

func (this *EvidenceController) UpdateEvidence(ctx echo.Context) error {

	evidenceId := ctx.Param("evidence-id")
	req := new(evd_req.UpdateEvdReq)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, global_response.StandartResponse{
			Message: "Request not valid.",
		})
	}
	if err := ctx.Validate(req); err != nil {
		return err
	}

	userId, _ := lib.ExtractToken(ctx)
	evidence := dto.Evidence{
		Content: req.Content,
	}

	err := this.evdSvc.UpdateEvidence(userId, evidenceId, evidence)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when update evidence."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Evidence not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success update evidence.",
	})

}

func (this *EvidenceController) UpdateImage(ctx echo.Context) error {

	evidenceId := ctx.Param("evidence-id")

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

	userId, _ := lib.ExtractToken(ctx)

	err = this.evdSvc.UpdateImage(userId, evidenceId, file.Filename, src)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when update image."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Evidence not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success update image.",
	})

}

func (this *EvidenceController) DeleteEvidence(ctx echo.Context) error {

	evidenceId := ctx.Param("evidence-id")
	userId, _ := lib.ExtractToken(ctx)

	err := this.evdSvc.DeleteEvidence(userId, evidenceId)
	if err != nil {

		fmt.Println(err.Error())
		errMessage := err.Error()
		errResMessage := "Error when delete evidence."
		errResStatus := http.StatusInternalServerError

		if errMessage == "record not found" {
			errResMessage = "Evidence not found."
			errResStatus = http.StatusNotFound
		}

		return ctx.JSON(errResStatus, global_response.StandartResponse{
			Message: errResMessage,
		})

	}

	return ctx.JSON(http.StatusOK, global_response.StandartResponse{
		Message: "Success delete evidence.",
	})

}

package lib

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ReqError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {

	var errors []*ReqError

	if err := cv.Validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el ReqError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, &el)
		}

		return echo.NewHTTPError(http.StatusBadRequest, errors)
	}

	return nil
}

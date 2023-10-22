package validator

import (
	"errors"
	"mime/multipart"
	"miniproject-alterra/app/lib"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ImageValidation(file *multipart.FileHeader) bool {

	extension := strings.ToLower(filepath.Ext(file.Filename))
	validImageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	if !lib.Contains(validImageExtensions, extension) {
		return false
	}

	contentType := file.Header.Get("Content-Type")
	validImageMimeTypes := []string{"image/jpeg", "image/png", "image/gif", "image/bmp"}
	if !lib.Contains(validImageMimeTypes, contentType) {
		return false
	}

	return true

}

func DateValidation(err error, dob time.Time) error {

	if err != nil {
		return err
	}

	minConstraint := time.Date(1945, 8, 17, 0, 0, 0, 0, time.UTC)
	if dob.After(time.Now()) || dob.Before(minConstraint) {
		return errors.New("Time threeshold exceeded.")
	}

	return nil

}

func GoogleMapsURLValidator(str string) bool {

	var reg1 *regexp.Regexp
	var reg2 *regexp.Regexp
	var err error

	reg1, err = regexp.Compile(`^https://www\.google\.co\.id/maps/`)
	if err != nil {
		return false
	}
	reg2, err = regexp.Compile(`^https://maps\.app\.goo\.gl/`)
	if err != nil {
		return false
	}

	if !reg1.MatchString(str) && !reg2.MatchString(str) {
		return false
	}

	return false

}

package validator

import (
	"errors"
	"mime/multipart"
	"miniproject-alterra/app/lib"
	"path/filepath"
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

	minConstraint := time.Date(1945, 8, 15, 0, 0, 0, 0, time.UTC)
	if dob.After(time.Now()) || dob.Before(minConstraint) {
		return errors.New("Time threeshold exceeded.")
	}

	return nil

}

package validator

import (
	"errors"
	"strings"
	"time"
)

var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func DetectImageType(b []byte) string {

	s := string(b)
	for key, val := range magicTable {
		if strings.HasPrefix(s, key) {
			return val
		}
	}

	return ""

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

package global_entity

import (
	"io"
	"miniproject-alterra/module/dto"
)

type EmailDataFormat struct {
	Name string
	URL  string
}

type SendEmailFormat struct {
	To      string
	Cc      string
	Subject string
	Body    string
}

type (
	EmailServiceInterface interface {
		SendEmail(sendEmailFormat SendEmailFormat) error
	}

	StorageServiceInterface interface {
		UploadFile(bucketName string, fileName string, body io.Reader) error
		DeleteFile(bucketName string, fileName string) error
		GetUrl(bucketName string, fileName string) (string, error)
	}

	IGlobalRepository interface {
		GetUser(userId string) (dto.User, error)
	}
)

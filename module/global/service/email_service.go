package global_service

import (
	"crypto/tls"
	"miniproject-alterra/app/config"
	global_entity "miniproject-alterra/module/global/entity"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.AppConfig
}

func NewEmailService(config *config.AppConfig) global_entity.EmailServiceInterface {

	return &EmailService{
		config: config,
	}

}

func (this *EmailService) SendEmail(format global_entity.SendEmailFormat) error {

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", this.config.EMAIL_SENDER_NAME)
	mailer.SetHeader("To", format.To)
	mailer.SetAddressHeader("Cc", format.Cc, "Test")
	mailer.SetHeader("Subject", format.Subject)
	mailer.SetBody("text/html", format.Body)

	dialer := gomail.NewDialer(
		this.config.EMAIL_SMTP_HOST,
		this.config.EMAIL_SMTP_PORT,
		this.config.EMAIL_SMTP_EMAIL,
		this.config.EMAIL_SMTP_PASSWORD,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil

}

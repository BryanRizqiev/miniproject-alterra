package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var JWT_SECRET string

type AppConfig struct {
	EMAIL_SENDER_NAME     string
	EMAIL_SMTP_HOST       string
	EMAIL_SMTP_PORT       int
	EMAIL_SMTP_EMAIL      string
	EMAIL_SMTP_PASSWORD   string
	DB_USERNAME           string
	DB_PASSWORD           string
	DB_HOSTNAME           string
	DB_PORT               string
	DB_NAME               string
	JWT_KEY               string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	ENDPOINT              string
}

func GetConfig() *AppConfig {

	return Config()

}

func Config() *AppConfig {

	var app AppConfig
	err := godotenv.Load()
	if err != nil {
		panic("Error when loading .env file")
	}

	app.DB_USERNAME = os.Getenv("DB_USERNAME")
	app.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	app.DB_HOSTNAME = os.Getenv("DB_HOSTNAME")
	app.DB_PORT = os.Getenv("DB_PORT")
	app.DB_NAME = os.Getenv("DB_NAME")
	app.JWT_KEY = os.Getenv("JWT_KEY")
	JWT_SECRET = app.JWT_KEY

	app.EMAIL_SENDER_NAME = os.Getenv("EMAIL_SENDER_NAME")
	app.EMAIL_SMTP_HOST = os.Getenv("EMAIL_SMTP_HOST")
	smptpPort, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		panic("SMTP port must number")
	}
	app.EMAIL_SMTP_PORT = smptpPort
	app.EMAIL_SMTP_EMAIL = os.Getenv("EMAIL_SMTP_EMAIL")
	app.EMAIL_SMTP_PASSWORD = os.Getenv("EMAIL_SMTP_PASSWORD")

	app.AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
	app.AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	app.AWS_REGION = os.Getenv("AWS_REGION")
	app.ENDPOINT = os.Getenv("ENDPOINT")

	return &app

}

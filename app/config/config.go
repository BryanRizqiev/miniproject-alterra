package config

var JWT_SECRET string

type AppConfig struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOSTNAME string
	DB_PORT     string
	DB_NAME     string
	JWT_KEY     string
}

func GetConfig() *AppConfig {

	return Config()

}

func Config() *AppConfig {

	var app AppConfig
	// app.DB_USERNAME = os.Getenv("DB_USERNAME")
	// app.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	// app.DB_HOSTNAME = os.Getenv("DB_HOSTNAME")
	// app.DB_PORT = os.Getenv("DB_PORT")
	// app.DB_NAME = os.Getenv("DB_NAME")
	// app.JWT_KEY = os.Getenv("JWT_KEY")
	app.DB_USERNAME = "root"
	app.DB_PASSWORD = ""
	app.DB_HOSTNAME = "127.0.0.1"
	app.DB_PORT = "3306"
	app.DB_NAME = "miniproject"
	app.JWT_KEY = "TOP_SECRET"
	JWT_SECRET = app.JWT_KEY

	return &app

}

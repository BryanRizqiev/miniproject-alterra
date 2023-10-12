package database

import (
	"fmt"
	"miniproject-alterra/app/config"
	user_model "miniproject-alterra/module/user/repository/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDBMysql(cfg *config.AppConfig) *gorm.DB {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOSTNAME, cfg.DB_PORT, cfg.DB_NAME)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db

}

func InitMigrationMySQL(db *gorm.DB) error {

	db.Migrator().DropTable(user_model.User{})

	return db.AutoMigrate(
		user_model.User{},
	)

}

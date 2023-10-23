package global_repo

import (
	"miniproject-alterra/module/dto"
	global_entity "miniproject-alterra/module/global/entity"

	"gorm.io/gorm"
)

type GlobalRepo struct {
	db *gorm.DB
}

func NewGlobalRepo(db *gorm.DB) global_entity.IGlobalRepository {

	return &GlobalRepo{
		db: db,
	}

}

func (this *GlobalRepo) GetUser(userId string) (dto.User, error) {

	var user dto.User

	tx := this.db.First(&user, "id = ?", userId)
	if tx.Error != nil {
		return dto.User{}, tx.Error
	}

	return user, nil

}

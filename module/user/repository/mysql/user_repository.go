package mysql_user_repository

import (
	"miniproject-alterra/app/lib"
	user_entity "miniproject-alterra/module/user/entity"
	user_model "miniproject-alterra/module/user/repository/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user_entity.UserRepositoryInterface {

	return &UserRepository{
		db: db,
	}

}

func (*UserRepository) CheckUser(userDTO user_entity.UserDTO) (UserDTO error) {

	panic("unimplemented")

}

func (*UserRepository) GetAllUser() ([]user_entity.UserDTO, error) {

	panic("unimplemented")

}

func (this *UserRepository) InsertUser(userDTO user_entity.UserDTO) error {

	user := user_model.User{
		Email:    userDTO.Email,
		Name:     userDTO.Name,
		Password: userDTO.Password,
		Address:  lib.NewNullString(userDTO.Address),
	}

	if tx := this.db.Create(&user); tx.Error != nil {
		return tx.Error
	}

	return nil

}

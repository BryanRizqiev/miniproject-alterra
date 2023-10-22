package mysql_user_repository

import (
	"miniproject-alterra/app/lib"
	user_entity "miniproject-alterra/module/user/entity"
	user_model "miniproject-alterra/module/user/repository/model"
	"time"

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

func (this *UserRepository) GetUserByEmail(email string) (user_entity.UserDTO, error) {

	var user user_model.User
	tx := this.db.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return user_entity.UserDTO{}, tx.Error
	}

	userDTO := user_entity.UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}

	return userDTO, nil

}

func (*UserRepository) GetAllUser() ([]user_entity.UserDTO, error) {

	panic("unimplemented")

}

func (this *UserRepository) InsertUser(userDTO user_entity.UserDTO) error {

	testTime := time.Time{}

	user := user_model.User{
		ID:              lib.NewUuid(),
		Name:            userDTO.Name,
		Email:           userDTO.Email,
		Password:        userDTO.Password,
		DOB:             lib.NewNullString(userDTO.DOB),
		Phone:           lib.NewNullString(userDTO.Phone),
		Address:         lib.NewNullString(userDTO.Address),
		VerifiedEmailAt: lib.NewNullTime(testTime),
	}

	tx := this.db.Create(&user)

	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

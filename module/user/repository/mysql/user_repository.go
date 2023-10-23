package mysql_user_repository

import (
	"errors"
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

func (this *UserRepository) UpdateUserRequestVerified(userId string) error {

	var user user_model.User

	tx := this.db.First(&user, "id = ?", userId)
	if tx.Error != nil {
		return tx.Error
	}
	tx = this.db.Model(&user).Update("request_verified", "request")
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (this *UserRepository) UpdateUserVerifiedEmail(userId string) error {

	var user user_model.User

	tx := this.db.First(&user, "id = ?", userId)
	if tx.Error != nil {
		return tx.Error
	}

	if !user.VerifiedEmailAt.Time.IsZero() {
		return errors.New("user already verified")
	}
	tx = this.db.Model(&user).Update("verified_email_at", time.Now())
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (this *UserRepository) CheckUserVerifiedEmail(userId string) (bool, error) {

	var user user_model.User

	tx := this.db.First(&user, "id = ?", userId)
	if tx.Error != nil {
		return false, tx.Error
	}

	if user.VerifiedEmailAt.Time.IsZero() {
		return false, nil
	}

	return true, nil

}

func (this *UserRepository) UpdateUserRole(userId string, role string) error {

	var user user_model.User

	tx := this.db.First(&user, "id = ?", userId)
	if tx.Error != nil {
		return tx.Error
	}

	tx = this.db.Model(&user).Update("role", role)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (this *UserRepository) FindUser(userId string) (user_model.User, error) {

	var user user_model.User

	tx := this.db.First(&user, "id = ?", userId)
	if tx.Error != nil {
		return user_model.User{}, tx.Error
	}

	return user, nil

}

package mysql_user_repository

import (
	"errors"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
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

func (this *UserRepository) GetUserByEmail(email string) (dto.User, error) {

	var user dto.User
	tx := this.db.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return dto.User{}, tx.Error
	}

	return user, nil

}

func (this *UserRepository) GetAllUser() ([]dto.User, error) {

	panic("unimplemented")

}

func (this *UserRepository) InsertUser(user dto.User) error {

	userInsert := user_model.User{
		ID:       lib.NewUuid(),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		DOB:      lib.NewNullString(user.DOB.String),
		Phone:    lib.NewNullString(user.Phone.String),
		Address:  lib.NewNullString(user.Address.String),
	}

	tx := this.db.Create(&userInsert)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (this *UserRepository) UpdateUserRequestVerified(userId string) error {

	var user dto.User

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

	var user dto.User

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

	var user dto.User

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

	var user dto.User

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

func (this *UserRepository) FindUser(userId string) (dto.User, error) {

	var user dto.User

	tx := this.db.First(&user, "id = ?", userId)
	if tx.Error != nil {
		return dto.User{}, tx.Error
	}

	return user, nil

}

func (this *UserRepository) GetRequestingUser() ([]dto.User, error) {

	var users []dto.User

	err := this.db.Where("request_verified = ?", "request").Find(&users).Error
	if err != nil {
		return []dto.User{}, err
	}

	return users, nil

}

package user_service

import (
	"errors"
	"fmt"
	"miniproject-alterra/app/config"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
	global_entity "miniproject-alterra/module/global/entity"
	user_entity "miniproject-alterra/module/user/entity"
)

type UserService struct {
	userRepo     user_entity.UserRepositoryInterface
	emailService global_entity.EmailServiceInterface
	config       *config.AppConfig
}

func NewUserService(userRepo user_entity.UserRepositoryInterface, emailService global_entity.EmailServiceInterface, config *config.AppConfig) user_entity.UserServiceInterface {

	return &UserService{
		userRepo:     userRepo,
		emailService: emailService,
		config:       config,
	}

}

func (this *UserService) Login(user dto.User) (string, error) {

	var err error
	var token string
	oldPassword := user.Password

	user, err = this.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}
	if !lib.BcryptMatchingPassword(user.Password, oldPassword) {
		return "", errors.New("credentials not valid")
	}

	token, err = lib.CreateToken(user.Id, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (this *UserService) Register(user dto.User) error {

	encryptedPassword, _ := lib.BcryptHashPassword(user.Password)

	newUser := dto.User{
		Email:    user.Email,
		Name:     user.Name,
		Password: encryptedPassword,
		DOB:      user.DOB,
		Address:  user.Address,
		Phone:    user.Phone,
	}

	if err := this.userRepo.InsertUser(newUser); err != nil {
		return err
	}

	return nil

}

func (this *UserService) RequestVerified(userId string) error {

	isVerifiedEmail, err := this.userRepo.CheckUserVerifiedEmail(userId)
	if err != nil {
		return err
	}
	if !isVerifiedEmail {
		return errors.New("email not verified")
	}

	err = this.userRepo.UpdateUserRequestVerified(userId)
	if err != nil {
		return err
	}

	return nil

}

func (this *UserService) VerifyEmail(userId string) error {

	err := this.userRepo.UpdateUserVerifiedEmail(userId)
	if err != nil {
		return err
	}

	return nil
}

func (this *UserService) RequestVerifyEmail(userId string) error {

	user, err := this.userRepo.FindUser(userId)
	if err != nil {
		return err
	}
	if !user.VerifiedEmailAt.Time.IsZero() {
		return errors.New("user already verified")
	}

	emailData := global_entity.EmailDataFormat{
		Name: user.Name,
		URL:  fmt.Sprintf("http://%s/verify-email/%s", this.config.APP_URL, user.Id),
	}
	htmlStr, err := lib.ParseTemplate("./app/lib/template/email.html", emailData)
	if err != nil {
		return err
	}
	format := global_entity.SendEmailFormat{
		To:      user.Email,
		Cc:      user.Email,
		Subject: "Email Verification",
		Body:    htmlStr,
	}
	go this.emailService.SendEmail(format)

	return nil

}

func (this *UserService) GetRequestingUser(userId string) ([]dto.User, error) {

	user, err := this.userRepo.FindUser(userId)
	if err != nil {
		return []dto.User{}, err
	}

	if !lib.CheckIsAdmin(user) {
		return []dto.User{}, errors.New("user not allowed")
	}

	users, err := this.userRepo.GetRequestingUser()
	if err != nil {
		return []dto.User{}, err
	}

	return users, nil

}

func (this *UserService) ChangeUserRole(reqUserId string, userId string, role string) error {

	user, err := this.userRepo.FindUser(reqUserId)
	if err != nil {
		return err
	}

	if !lib.CheckIsAdmin(user) {
		return errors.New("user not allowed")
	}

	err = this.userRepo.UpdateUserRole(userId, role)
	if err != nil {
		return err
	}

	return nil

}

func (this *UserService) GetAllUser(userId string) ([]dto.User, error) {

	user, err := this.userRepo.FindUser(userId)
	if err != nil {
		return []dto.User{}, nil
	}

	if !lib.CheckIsAdmin(user) {
		return []dto.User{}, errors.New("user not allowed")
	}

	users, err := this.userRepo.GetAllUser()
	if err != nil {
		return []dto.User{}, nil
	}

	return users, nil

}

func (this *UserService) UpdateUser(userId string, payload dto.User) error {

	user, err := this.userRepo.FindUser(userId)
	if err != nil {
		return err
	}

	user.Name = payload.Name
	user.Address = payload.Address
	user.DOB = payload.DOB
	user.Phone = payload.Phone

	err = this.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil

}

func (this *UserService) DeleteUser(reqUserId, userId string) error {

	user, err := this.userRepo.FindUser(reqUserId)
	if err != nil {
		return err
	}
	if !lib.CheckIsAdmin(user) {
		return errors.New("user not allowed")
	}

	user, err = this.userRepo.FindUser(userId)
	if err != nil {
		return err
	}
	err = this.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}

	return nil

}

func (this *UserService) UserSelfDelete(userId string) error {

	user, err := this.userRepo.FindUser(userId)
	if err != nil {
		return err
	}
	if lib.CheckIsAdmin(user) {
		return errors.New("user not allowed")
	}

	err = this.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}

	return nil

}

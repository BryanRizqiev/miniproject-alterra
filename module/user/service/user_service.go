package user_service

import (
	"errors"
	"fmt"
	"miniproject-alterra/app/config"
	"miniproject-alterra/app/lib"
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

func (this *UserService) Login(userDTO user_entity.UserDTO) (string, error) {

	var err error
	var token string
	oldPassword := userDTO.Password

	userDTO, err = this.userRepo.GetUserByEmail(userDTO.Email)
	if err != nil {
		return "", err
	}
	if !lib.BcryptMatchingPassword(userDTO.Password, oldPassword) {
		return "", errors.New("credentials not valid")
	}

	token, err = lib.CreateToken(userDTO.ID, userDTO.Email)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (this *UserService) GetAllUser() ([]user_entity.UserDTO, error) {

	panic("unimplemented")

}

func (this *UserService) Register(userDTO user_entity.UserDTO) error {

	encryptedPassword, _ := lib.BcryptHashPassword(userDTO.Password)

	newUserDTO := user_entity.UserDTO{
		Email:    userDTO.Email,
		Name:     userDTO.Name,
		Password: encryptedPassword,
		DOB:      userDTO.DOB,
		Address:  userDTO.Address,
		Phone:    userDTO.Phone,
	}

	if err := this.userRepo.InsertUser(newUserDTO); err != nil {
		return err
	}

	return nil

}

func (this *UserService) RequestVerification(userId string) error {

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

func (this *UserService) Verify(userId string) error {

	err := this.userRepo.UpdateUserVerifiedEmail(userId)
	if err != nil {
		return err
	}

	return nil
}

func (this *UserService) SendVerifyEmail(userId string) error {

	user, err := this.userRepo.FindUser(userId)
	if err != nil {
		return err
	}

	if !user.VerifiedEmailAt.Time.IsZero() {
		return errors.New("user already verified")
	}

	emailData := global_entity.EmailDataFormat{
		Name: user.Name,
		URL:  fmt.Sprintf("http://%s/verify-email/%s", this.config.APP_URL, user.ID),
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

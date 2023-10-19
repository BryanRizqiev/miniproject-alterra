package user_service

import (
	"errors"
	"miniproject-alterra/app/lib"
	user_entity "miniproject-alterra/module/user/entity"
)

type UserService struct {
	userRepository user_entity.UserRepositoryInterface
}

func NewUserService(userRepository user_entity.UserRepositoryInterface) user_entity.UserServiceInterface {

	return &UserService{
		userRepository: userRepository,
	}

}

func (this *UserService) Login(userDTO user_entity.UserDTO) (string, error) {

	var err error
	var token string
	oldPassword := userDTO.Password

	userDTO, err = this.userRepository.GetUserByEmail(userDTO.Email)
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

	if err := this.userRepository.InsertUser(newUserDTO); err != nil {
		return err
	}

	return nil

}

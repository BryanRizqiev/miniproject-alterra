package user_service_test

import (
	"errors"
	"miniproject-alterra/app/config"
	"miniproject-alterra/mocks"
	"miniproject-alterra/module/dto"
	user_service "miniproject-alterra/module/user/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {

	cfg := &config.AppConfig{}
	userRepo := new(mocks.UserRepository)
	storageSvc := new(mocks.StorageService)
	emailSvc := new(mocks.EmailService)
	expectedErr := errors.New("an error ocurred")

	var testCases = []struct {
		name string
		err  error
		user dto.User
	}{
		{
			name: "Error register when repository return error",
			err:  expectedErr,
			user: dto.User{
				Password: "any",
			},
		},
		{
			name: "Success register",
			err:  nil,
			user: dto.User{
				Password: "any",
			},
		},
	}

	for idx, testCase := range testCases {

		if idx == 0 {
			userRepo.On("InsertUser", mock.Anything).Return(expectedErr).Once()
		}
		if idx == 1 {
			userRepo.On("InsertUser", mock.Anything).Return(nil).Once()
		}

		t.Run(testCase.name, func(t *testing.T) {
			svc := user_service.NewUserService(userRepo, emailSvc, storageSvc, cfg)

			err := svc.Register(testCase.user)
			if idx == 1 {
				assert.NoError(t, err)
				assert.Equal(t, testCase.err, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, testCase.err, err)
			}
		})

		t.Cleanup(func() {
			userRepo.AssertExpectations(t)
			storageSvc.AssertExpectations(t)
			emailSvc.AssertExpectations(t)
		})

	}

}

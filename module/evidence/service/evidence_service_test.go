package evd_svc_test

import (
	"errors"
	"mime/multipart"
	"miniproject-alterra/mocks"
	"miniproject-alterra/module/dto"
	evd_svc "miniproject-alterra/module/evidence/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEvidence(t *testing.T) {

	evidenceRepo := new(mocks.EvidenceRepository)
	storageSvc := new(mocks.StorageService)
	globalRepo := new(mocks.GlobalRepo)
	expectedErr := errors.New("an error ocurred")

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "Error create evidence when repository return error",
			err:  expectedErr,
		},
		{
			name: "Error create evidence when storage service return error",
			err:  expectedErr,
		},
		{
			name: "Success create evidence",
			err:  nil,
		},
	}

	storageSvc.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(expectedErr).Once()
	evidenceRepo.On("InsertEvidence", mock.Anything).Return(expectedErr).Once()

	for idx, testCase := range testCases {

		if idx == 1 {
			storageSvc.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(nil).Twice()
		}
		if idx == 2 {
			evidenceRepo.On("InsertEvidence", mock.Anything).Return(nil).Once()
		}

		t.Run(testCase.name, func(t *testing.T) {
			svc := evd_svc.NewEvidenceService(evidenceRepo, storageSvc, globalRepo)

			var image multipart.File
			err := svc.CreateEvidence("any", "any", image, dto.Evidence{})
			if idx == 2 {
				assert.NoError(t, err)
				assert.Equal(t, testCase.err, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, testCase.err, err)
			}
		})

		t.Cleanup(func() {
			globalRepo.AssertExpectations(t)
			evidenceRepo.AssertExpectations(t)
			storageSvc.AssertExpectations(t)
		})

	}

}

func TestGetEvidences(t *testing.T) {

	evidenceRepo := new(mocks.EvidenceRepository)
	storageSvc := new(mocks.StorageService)
	globalRepo := new(mocks.GlobalRepo)
	expectedErr := errors.New("an error ocurred")

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "Error get evidence when repository return error-1",
			err:  expectedErr,
		},
		{
			name: "Error get evidence when check is admin",
			err:  errors.New("user not allowed"),
		},
		{
			name: "Error get evidence when repository return error-2",
			err:  expectedErr,
		},
		{
			name: "Error get evidence when storage service return error",
			err:  expectedErr,
		},
		{
			name: "Success get evidences",
			err:  nil,
		},
	}

	user := dto.User{
		Role: "user",
	}
	admin := dto.User{
		Role: "admin",
	}
	evidences := []dto.Evidence{
		{
			Id: "1",
		},
	}

	for idx, testCase := range testCases {

		if idx == 0 {
			globalRepo.On("GetUser", mock.AnythingOfType("string")).Return(dto.User{}, expectedErr).Once()
		}
		if idx == 1 {
			globalRepo.On("GetUser", mock.AnythingOfType("string")).Return(user, nil).Once()
		}
		if idx == 2 {
			globalRepo.On("GetUser", mock.AnythingOfType("string")).Return(admin, nil).Times(3)
			evidenceRepo.On("GetEvidences", mock.AnythingOfType("string")).Return([]dto.Evidence{}, expectedErr).Once()
		}
		if idx == 3 {
			evidenceRepo.On("GetEvidences", mock.AnythingOfType("string")).Return(evidences, nil).Twice()
			storageSvc.On("GetUrl", mock.AnythingOfType("string"), mock.Anything).Return("", expectedErr).Once()
		}
		if idx == 4 {
			storageSvc.On("GetUrl", mock.AnythingOfType("string"), mock.Anything).Return("the-url", nil).Once()
		}

		t.Run(testCase.name, func(t *testing.T) {
			svc := evd_svc.NewEvidenceService(evidenceRepo, storageSvc, globalRepo)
			_, err := svc.GetEvidences("any", "any")
			if idx == 4 {
				assert.NoError(t, err)
				assert.Equal(t, testCase.err, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, testCase.err, err)
			}
		})

		t.Cleanup(func() {
			globalRepo.AssertExpectations(t)
			evidenceRepo.AssertExpectations(t)
			storageSvc.AssertExpectations(t)
		})

	}

}

func TestUpdateEvidences(t *testing.T) {

	evidenceRepo := new(mocks.EvidenceRepository)
	storageSvc := new(mocks.StorageService)
	globalRepo := new(mocks.GlobalRepo)
	expectedErr := errors.New("an error ocurred")

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "Error update evidence when repository return error-1",
			err:  expectedErr,
		},
		{
			name: "Error update evidence when repository return error-2",
			err:  expectedErr,
		},
		{
			name: "Success update evidences",
			err:  nil,
		},
	}

	for idx, testCase := range testCases {

		if idx == 0 {
			evidenceRepo.On("FindOwnEvidence", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dto.Evidence{}, expectedErr).Once()
		}
		if idx == 1 {
			evidenceRepo.On("FindOwnEvidence", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dto.Evidence{}, nil).Twice()
			evidenceRepo.On("UpdateEvidence", mock.Anything).Return(expectedErr).Once()
		}
		if idx == 1 {
			evidenceRepo.On("UpdateEvidence", mock.Anything).Return(nil).Once()
		}

		t.Run(testCase.name, func(t *testing.T) {
			svc := evd_svc.NewEvidenceService(evidenceRepo, storageSvc, globalRepo)
			err := svc.UpdateEvidence("any", "any", dto.Evidence{})
			if idx == 2 {
				assert.NoError(t, err)
				assert.Equal(t, testCase.err, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, testCase.err, err)
			}
		})

		t.Cleanup(func() {
			globalRepo.AssertExpectations(t)
			evidenceRepo.AssertExpectations(t)
			storageSvc.AssertExpectations(t)
		})

	}

}

func TestUpdateImage(t *testing.T) {

	evidenceRepo := new(mocks.EvidenceRepository)
	storageSvc := new(mocks.StorageService)
	globalRepo := new(mocks.GlobalRepo)
	expectedErr := errors.New("an error ocurred")

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "Error update image when repository return error-1",
			err:  expectedErr,
		},
		{
			name: "Error update image when storage service return error",
			err:  expectedErr,
		},
		{
			name: "Error update image when repository return error-2",
			err:  expectedErr,
		},
		{
			name: "Success update image",
			err:  nil,
		},
	}

	for idx, testCase := range testCases {

		if idx == 0 {
			evidenceRepo.On("FindOwnEvidence", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dto.Evidence{}, expectedErr).Once()
		}
		if idx == 1 {
			evidenceRepo.On("FindOwnEvidence", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dto.Evidence{}, nil).Times(3)
			storageSvc.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(expectedErr).Once()
		}
		if idx == 2 {
			storageSvc.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			evidenceRepo.On("UpdateImage", mock.AnythingOfType("string"), mock.Anything).Return(expectedErr).Once()
		}
		if idx == 3 {
			evidenceRepo.On("UpdateImage", mock.AnythingOfType("string"), mock.Anything).Return(nil).Once()
			storageSvc.On("DeleteFile", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
		}

		t.Run(testCase.name, func(t *testing.T) {
			svc := evd_svc.NewEvidenceService(evidenceRepo, storageSvc, globalRepo)

			var image multipart.File
			err := svc.UpdateImage("any", "any", "any", image)
			if idx == 3 {
				assert.NoError(t, err)
				assert.Equal(t, testCase.err, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, testCase.err, err)
			}
		})

		t.Cleanup(func() {
			globalRepo.AssertExpectations(t)
			evidenceRepo.AssertExpectations(t)
			storageSvc.AssertExpectations(t)
		})

	}

}

func TestDeleteEvidence(t *testing.T) {

	evidenceRepo := new(mocks.EvidenceRepository)
	storageSvc := new(mocks.StorageService)
	globalRepo := new(mocks.GlobalRepo)
	expectedErr := errors.New("an error ocurred")

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "Error delete evidence when repository return error-1",
			err:  expectedErr,
		},
		{
			name: "Error delete evidence when repository return error-2",
			err:  expectedErr,
		},
		{
			name: "Error delete evidence when repository return error-3",
			err:  expectedErr,
		},
		{
			name: "Success delete evidence",
			err:  nil,
		},
	}

	admin := dto.User{
		Role: "admin",
	}

	for idx, testCase := range testCases {

		if idx == 0 {
			globalRepo.On("GetUser", mock.AnythingOfType("string")).Return(dto.User{}, expectedErr).Once()
		}
		if idx == 1 {
			globalRepo.On("GetUser", mock.AnythingOfType("string")).Return(admin, nil).Times(3)
			evidenceRepo.On("FindEvidence", mock.AnythingOfType("string")).Return(dto.Evidence{}, expectedErr).Once()
		}
		if idx == 2 {
			evidenceRepo.On("FindEvidence", mock.AnythingOfType("string")).Return(dto.Evidence{}, nil).Twice()
			evidenceRepo.On("DeleteEvidence", mock.Anything).Return(expectedErr).Once()
		}
		if idx == 3 {
			evidenceRepo.On("DeleteEvidence", mock.Anything).Return(nil).Once()
		}

		t.Run(testCase.name, func(t *testing.T) {
			svc := evd_svc.NewEvidenceService(evidenceRepo, storageSvc, globalRepo)

			err := svc.DeleteEvidence("any", "any")
			if idx == 3 {
				assert.NoError(t, err)
				assert.Equal(t, testCase.err, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, testCase.err, err)
			}
		})

		t.Cleanup(func() {
			globalRepo.AssertExpectations(t)
			evidenceRepo.AssertExpectations(t)
			storageSvc.AssertExpectations(t)
		})

	}

}

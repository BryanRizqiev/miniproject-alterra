package evd_svc

import (
	"errors"
	"fmt"
	"mime/multipart"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
	evd_entity "miniproject-alterra/module/evidence/entity"
	global_entity "miniproject-alterra/module/global/entity"
	"path/filepath"
	"strings"
)

type EvidenceService struct {
	evidenceRepo evd_entity.IEvidenceRepository
	storageSvc   global_entity.StorageServiceInterface
	globalRepo   global_entity.IGlobalRepository
}

func NewEvidenceService(
	evidenceRepo evd_entity.IEvidenceRepository,
	storageSvc global_entity.StorageServiceInterface,
	globalRepo global_entity.IGlobalRepository,
) evd_entity.IEvidenceService {

	return &EvidenceService{
		evidenceRepo: evidenceRepo,
		storageSvc:   storageSvc,
		globalRepo:   globalRepo,
	}

}

func (this *EvidenceService) CreateEvidence(userId string, eventId string, image multipart.File, evidence dto.Evidence) error {

	evidence.CreatedBy = userId
	evidence.EventId = eventId

	fileExt := strings.ToLower(filepath.Ext(evidence.Image))
	newFilename := fmt.Sprintf("%s-%s%s", "evidence", lib.RandomString(16), fileExt)
	evidence.Image = newFilename

	var err error
	err = this.storageSvc.UploadFile("event-evidence", newFilename, image)
	if err != nil {
		return err
	}

	err = this.evidenceRepo.InsertEvidence(evidence)
	if err != nil {
		return err
	}

	return nil

}

func (this *EvidenceService) GetEvidences(userId, eventId string) ([]dto.Evidence, error) {

	user, err := this.globalRepo.GetUser(userId)
	if err != nil {
		return []dto.Evidence{}, err
	}

	if !lib.CheckIsAdmin(user) {
		return []dto.Evidence{}, errors.New("user not allowed")
	}

	evidences, err := this.evidenceRepo.GetEvidences(eventId)
	if err != nil {
		return []dto.Evidence{}, err
	}

	for i := range evidences {
		url, errURL := this.storageSvc.GetUrl("event-evidence", evidences[i].Image)
		if errURL != nil {
			return []dto.Evidence{}, err
		}
		evidences[i].Image = url
	}

	return evidences, nil

}

func (this *EvidenceService) UpdateEvidence(userId, evidenceId string, payload dto.Evidence) error {

	evidence, err := this.evidenceRepo.FindOwnEvidence(userId, evidenceId)
	if err != nil {
		return err
	}

	evidence.Content = payload.Content

	err = this.evidenceRepo.UpdateEvidence(evidence)
	if err != nil {
		return err
	}

	return nil

}

func (this *EvidenceService) UpdateImage(userId, evidenceId, filename string, image multipart.File) error {

	evidence, err := this.evidenceRepo.FindOwnEvidence(userId, evidenceId)
	if err != nil {
		return err
	}

	fileExt := strings.ToLower(filepath.Ext(filename))
	newFilename := fmt.Sprintf("%s-%s%s", "evidence", lib.RandomString(16), fileExt)

	err = this.storageSvc.UploadFile("event-evidence", newFilename, image)
	if err != nil {
		return err
	}
	err = this.evidenceRepo.UpdateImage(newFilename, evidence)
	if err != nil {
		return err
	}
	_ = this.storageSvc.DeleteFile("event-evidence", evidence.Image)

	return nil

}

func (this *EvidenceService) DeleteEvidence(userId, evidenceId string) error {

	user, err := this.globalRepo.GetUser(userId)
	if err != nil {
		return err
	}
	if lib.CheckIsAdmin(user) {
		evidence, err := this.evidenceRepo.FindEvidence(evidenceId)
		if err != nil {
			return err
		}

		err = this.evidenceRepo.DeleteEvidence(evidence)
		if err != nil {
			return err
		}
	} else {
		evidence, err := this.evidenceRepo.FindOwnEvidence(userId, evidenceId)
		if err != nil {
			return err
		}

		err = this.evidenceRepo.DeleteEvidence(evidence)
		if err != nil {
			return err
		}
	}

	return nil

}

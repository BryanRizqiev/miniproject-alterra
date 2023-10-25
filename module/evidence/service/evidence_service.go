package evd_svc

import (
	"fmt"
	"mime/multipart"
	"miniproject-alterra/app/lib"
	evd_entity "miniproject-alterra/module/evidence/entity"
	global_entity "miniproject-alterra/module/global/entity"
	"path/filepath"
	"strings"
)

type EvidenceService struct {
	evdRepo evd_entity.IEvidenceRepository
	strgSvc global_entity.StorageServiceInterface
}

func NewEvidenceService(evdRepo evd_entity.IEvidenceRepository, strgSvc global_entity.StorageServiceInterface) evd_entity.IEvidenceService {

	return &EvidenceService{
		evdRepo: evdRepo,
		strgSvc: strgSvc,
	}

}

func (this *EvidenceService) CreateEvidence(userId string, evtId string, image multipart.File, evdD evd_entity.EvidenceDTO) error {

	evdD.UserId = userId
	evdD.EventId = evtId

	fileExt := strings.ToLower(filepath.Ext(evdD.Image))
	newFilename := fmt.Sprintf("%s-%s%s", "evidence", lib.RandomString(8), fileExt)
	evdD.Image = newFilename

	var err error
	err = this.strgSvc.UploadFile("event-evidence", newFilename, image)
	if err != nil {
		return err
	}

	err = this.evdRepo.InsertEvidence(evdD)
	if err != nil {
		return err
	}

	return nil

}

func (this *EvidenceService) GetEvidences(eventId string) ([]evd_entity.EvidenceDTO, error) {

	evdsD, err := this.evdRepo.GetEvidences(eventId)
	if err != nil {
		return []evd_entity.EvidenceDTO{}, err
	}

	for i := range evdsD {
		url, errURL := this.strgSvc.GetUrl("event", evdsD[i].Image)
		if errURL != nil {
			return []evd_entity.EvidenceDTO{}, err
		}
		evdsD[i].Image = url
	}

	return evdsD, nil

}

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

	evdD.UserID = userId
	evdD.EventID = evtId

	fileExt := strings.ToLower(filepath.Ext(evdD.Image))
	newFilename := fmt.Sprintf("%s-%s%s", "event", lib.RandomString(8), fileExt)
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

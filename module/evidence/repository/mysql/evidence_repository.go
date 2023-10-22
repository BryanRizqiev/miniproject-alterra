package mysql_evd_repo

import (
	"miniproject-alterra/app/lib"
	evd_entity "miniproject-alterra/module/evidence/entity"
	evd_model "miniproject-alterra/module/evidence/repository/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EvidenceRepository struct {
	db *gorm.DB
}

func NewEvidenceRepository(db *gorm.DB) evd_entity.IEvidenceRepository {
	return &EvidenceRepository{
		db: db,
	}
}

func (this *EvidenceRepository) InsertEvidence(evdD evd_entity.EvidenceDTO) error {

	evd := evd_model.Evidence{
		ID:        lib.NewUuid(),
		Content:   evdD.Content,
		Image:     evdD.Image,
		CreatedBy: evdD.UserId,
		EventId:   evdD.EventId,
	}

	tx := this.db.Omit(clause.Associations).Create(&evd)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (this *EvidenceRepository) GetEvidences(eventId string) ([]evd_entity.EvidenceDTO, error) {

	var evdsM []evd_model.Evidence
	var evdsD []evd_entity.EvidenceDTO

	tx := this.db.Where("event_id = ?", eventId).Find(&evdsM)
	if tx.Error != nil {
		return evdsD, tx.Error
	}

	for _, value := range evdsM {
		evdD := evd_entity.EvidenceDTO{
			Id:      value.ID,
			Content: value.Content,
			Image:   value.Image,
		}
		evdsD = append(evdsD, evdD)
	}

	return evdsD, nil

}

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
		CreatedBy: evdD.UserID,
		EventId:   evdD.EventID,
	}

	tx := this.db.Omit(clause.Associations).Create(&evd)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

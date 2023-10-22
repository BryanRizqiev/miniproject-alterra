package mysql_evd_repo

import (
	evd_entity "miniproject-alterra/module/evidence/entity"

	"gorm.io/gorm"
)

type EvidenceRepository struct {
	db *gorm.DB
}

func NewEvidenceRepository(db *gorm.DB) evd_entity.IEvidenceRepository {
	return &EvidenceRepository{
		db: db,
	}
}

func (this *EvidenceRepository) Insert(evdD evd_entity.EvidenceDTO) error {
	panic("unimplemented")
}

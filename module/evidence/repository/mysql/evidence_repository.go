package mysql_evd_repo

import (
	"errors"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
	evd_entity "miniproject-alterra/module/evidence/entity"

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

func (this *EvidenceRepository) InsertEvidence(evidence dto.Evidence) error {

	evidence.Id = lib.NewUuid()

	tx := this.db.Omit(clause.Associations).Create(&evidence)
	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (this *EvidenceRepository) GetEvidences(eventId string) ([]dto.Evidence, error) {

	var evidences []dto.Evidence

	tx := this.db.Unscoped().Where("event_id = ?", eventId).Preload("User").Find(&evidences)
	if tx.Error != nil {
		return []dto.Evidence{}, tx.Error
	}

	return evidences, nil

}

func (this *EvidenceRepository) UpdateEvidence(evidence dto.Evidence) error {

	err := this.db.Save(&evidence).Error
	if err != nil {
		return err
	}

	return nil

}

func (this *EvidenceRepository) UpdateImage(fileName string, evidence dto.Evidence) error {

	err := this.db.Model(&evidence).Update("image", fileName).Error
	if err != nil {
		return err
	}

	return nil

}

func (this *EvidenceRepository) DeleteEvidence(event dto.Evidence) error {

	tx := this.db.Delete(&event)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected < 1 {
		return errors.New("nothing deleted")
	}

	return nil

}

func (this *EvidenceRepository) FindOwnEvidence(userId, evidenceId string) (dto.Evidence, error) {

	var evidence dto.Evidence

	tx := this.db.Where("created_by", userId).First(&evidence, "id = ?", evidenceId)
	if tx.Error != nil {
		return dto.Evidence{}, tx.Error
	}

	return evidence, nil

}

func (this *EvidenceRepository) FindEvidence(evidenceId string) (dto.Evidence, error) {

	var evidence dto.Evidence

	tx := this.db.First(&evidence, "id = ?", evidenceId)
	if tx.Error != nil {
		return dto.Evidence{}, tx.Error
	}

	return evidence, nil

}

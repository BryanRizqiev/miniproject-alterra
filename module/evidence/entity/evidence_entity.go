package evd_entity

import (
	"mime/multipart"
	"miniproject-alterra/module/dto"
)

type (
	IEvidenceService interface {
		CreateEvidence(userId string, eventId string, image multipart.File, evidence dto.Evidence) error
		GetEvidences(eventId string) ([]dto.Evidence, error)
		UpdateImage(userId, evidenceId, filename string, image multipart.File) error
		UpdateEvidence(userId, evidenceId string, payload dto.Evidence) error
		DeleteEvidence(userId, evidenceId string) error
	}

	IEvidenceRepository interface {
		InsertEvidence(evidence dto.Evidence) error
		GetEvidences(eventId string) ([]dto.Evidence, error)
		UpdateEvidence(evidence dto.Evidence) error
		DeleteEvidence(event dto.Evidence) error
		FindOwnEvidence(userId, evidenceId string) (dto.Evidence, error)
		UpdateImage(fileName string, evidence dto.Evidence) error
		FindEvidence(evidenceId string) (dto.Evidence, error)
	}
)

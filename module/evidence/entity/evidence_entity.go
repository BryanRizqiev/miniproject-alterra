package evd_entity

import (
	"mime/multipart"
	"miniproject-alterra/app/lib"
	"miniproject-alterra/module/dto"
	event_entity "miniproject-alterra/module/events/entity"
	"time"
)

type EvidenceDTO struct {
	Id        string
	Content   string
	Image     string
	UserId    string
	EventId   string
	Event     event_entity.EventDTO
	CreatedAt time.Time
}

type Evidence struct {
	*lib.Base

	Id        string
	Content   string
	Image     string
	CreatedBy string
	User      dto.User `gorm:"foreignKey:CreatedBy"`
	EventId   string
	Event     dto.Event `gorm:"foreignKey:EventId"`
}

type (
	IEvidenceService interface {
		CreateEvidence(userId string, evtId string, image multipart.File, evdD EvidenceDTO) error
		GetEvidences(eventId string) ([]EvidenceDTO, error)
	}

	IEvidenceRepository interface {
		InsertEvidence(evdD EvidenceDTO) error
		GetEvidences(eventId string) ([]EvidenceDTO, error)
	}
)

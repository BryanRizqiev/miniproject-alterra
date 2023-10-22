package evd_entity

import (
	"mime/multipart"
	event_entity "miniproject-alterra/module/events/entity"
	user_entity "miniproject-alterra/module/user/entity"
	"time"
)

type EvidenceDTO struct {
	Id        string
	Content   string
	Image     string
	UserId    string
	CreatedBy user_entity.UserDTO
	EventId   string
	Event     event_entity.EventDTO
	CreatedAt time.Time
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

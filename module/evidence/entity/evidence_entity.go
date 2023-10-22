package evd_entity

import (
	"mime/multipart"
	event_entity "miniproject-alterra/module/events/entity"
	user_entity "miniproject-alterra/module/user/entity"
	"time"
)

type EvidenceDTO struct {
	ID        string
	Content   string
	Image     string
	UserID    string
	CreatedBy user_entity.UserDTO
	EventID   string
	Event     event_entity.EventDTO
	CreatedAt time.Time
}

type (
	IEvidenceService interface {
		CreateEvidence(userId string, evtId string, image multipart.File, evdD EvidenceDTO) error
	}

	IEvidenceRepository interface {
		InsertEvidence(evdD EvidenceDTO) error
	}
)

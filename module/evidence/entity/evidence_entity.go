package evd_entity

import (
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
	EventId   string
	Event     event_entity.EventDTO
	CreatedAt time.Time
}

type (
	IEvidenceService interface {
		CreateEvidence(userId, evtId string, evdD EvidenceDTO) error
	}

	IEvidenceRepository interface {
		Insert(evdD EvidenceDTO) error
	}
)

package lib

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string `gorm:"type:string;primary_key;size:36;default:uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `grom:"index"`
}

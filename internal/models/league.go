package models

import (
	"github.com/google/uuid"
	"time"
)

type League struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Slug      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null"`
	Location  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

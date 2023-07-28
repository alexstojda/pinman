package models

import (
	"github.com/google/uuid"
	"time"
)

type League struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name string    `gorm:"type:varchar(255);not null"`

	// TODO: Limit to 20 characters and update request validation accordingly
	Slug       string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Owner      User
	OwnerID    uuid.UUID `gorm:"type:uuid;not null"`
	Location   Location
	LocationID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

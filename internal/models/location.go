package models

import (
	"github.com/google/uuid"
	"time"
)

type Location struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Slug         string    `gorm:"type:varchar(20);not null;uniqueIndex"`
	Address      string    `gorm:"type:varchar(255);not null"`
	PinballMapID int       `gorm:"type:int;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

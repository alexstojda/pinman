package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"pinman/internal/app/generated"
	"time"
)

type Tournament struct {
	ID         uuid.UUID                `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name       string                   `gorm:"type:varchar(255);not null"`
	Slug       string                   `gorm:"type:varchar(20);not null;uniqueIndex"`
	Type       generated.TournamentType `gorm:"type:varchar(100);not null"`
	Settings   datatypes.JSON           `gorm:"type:jsonb;not null"`
	LocationID uuid.UUID                `gorm:"type:uuid;not null"`
	Location   Location
	LeagueID   uuid.UUID `gorm:"type:uuid;not null"`
	League     League
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt  time.Time `gorm:"type:timestamp;not null;default:now()"`
}

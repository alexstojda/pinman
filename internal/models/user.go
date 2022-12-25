package models

import (
	"time"

	"github.com/google/uuid"
)

type UserClaims struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Role     string
	Verified bool
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"uniqueIndex"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"type:varchar(255);not null"`
	Verified  bool      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) GetUserClaims() UserClaims {
	return UserClaims{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Role:     u.Role,
		Verified: u.Verified,
	}
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null;index"`
	IsRevoked bool      `gorm:"not null;default:false;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

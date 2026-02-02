package model

import (
	"time"

	"github.com/google/uuid"
)

type MyPiece struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	Attendee   Attendee  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PieceCode  string    `gorm:"type:text;uniqueIndex;not null"`
	ExpireDate time.Time `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

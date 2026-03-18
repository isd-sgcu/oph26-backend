package entity

import (
	"time"

	"github.com/google/uuid"
)

type CollectedPiece struct {
	AttendeeID  uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	PieceID     uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	Attendee    Attendee  `gorm:"foreignKey:AttendeeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MyPiece     MyPiece   `gorm:"foreignKey:PieceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CollectedAt time.Time `gorm:"not null"`
}

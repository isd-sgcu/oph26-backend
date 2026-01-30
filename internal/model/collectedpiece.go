package model

import (
	"time"

	"github.com/google/uuid"
)

type CollectedPiece struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	PieceID     uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	User        User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MyPiece     MyPiece   `gorm:"foreignKey:PieceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CollectedAt time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

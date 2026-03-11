package entity

import (
	"time"

	"github.com/google/uuid"
)

type Checkin struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index:idx_user_faculty"`
	Faculty     string    `gorm:"not null;index:idx_user_faculty"`
	StaffID     uuid.UUID `gorm:"type:uuid;not null"`
	Staff       Staff     `gorm:"foreignKey:StaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CheckedInAt time.Time `gorm:"not null"`
}

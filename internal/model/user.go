package model

import (
	"time"

	"github.com/google/uuid"
)

// User model for database operations (GORM)
type User struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email      string     `gorm:"uniqueIndex;not null"`
	Role       string     `gorm:"not null"`
	AttendeeId *uuid.UUID `gorm:"type:uuid"`
	StaffId    *uuid.UUID `gorm:"type:uuid"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

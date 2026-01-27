package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email       string      `gorm:"uniqueIndex;not null"`
	Role        string      `gorm:"not null"`
	AttendeeId  *uuid.UUIDs `gorm:"type:uuid"`
	StaffId     *uuid.UUIDs `gorm:"type:uuid"`
	Score       *Score
	Leaderboard *Leaderboard
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

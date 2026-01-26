package model

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	// CUID is a unique identifier for student and staff members in Chula
	CUID      string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	FirstName string    `gorm:"type:varchar(100);not null"`
	LastName  string    `gorm:"type:varchar(100);not null"`
	Nickname  string    `gorm:"type:varchar(100)"`
	Phone     string    `gorm:"type:varchar(20)"`
	Year      string    `gorm:"type:varchar(10)"`
	Email     string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Faculty   string    `gorm:"type:varchar(100)"`
	// These are managed by GORM
	CreatedAt time.Time
	UpdatedAt time.Time
}

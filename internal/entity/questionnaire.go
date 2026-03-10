package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Questionnaire struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    string         `gorm:"type:text;not null;uniqueIndex"`
	User      User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Answers   datatypes.JSON `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
}

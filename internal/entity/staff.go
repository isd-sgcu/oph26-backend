package entity

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    *uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Cuid      string     `gorm:"type:text;uniqueIndex;not null"`
	Firstname string     `gorm:"type:text;not null"`
	Surname   string     `gorm:"type:text;not null"`
	Nickname  string     `gorm:"type:text;not null"`
	Phone     string     `gorm:"type:text;not null"`
	Year      string     `gorm:"type:text;not null"`
	Email     string     `gorm:"type:text;uniqueIndex;not null"`
	Faculty   string     `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

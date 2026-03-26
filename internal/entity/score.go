package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Score struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null;unique;primaryKey"`
	User    User
	Count1  int `gorm:"not null;default:0"`
	Count2  int `gorm:"not null;default:0"`
	Count3  int `gorm:"not null;default:0"`
	Count4  int `gorm:"not null;default:0"`
	Count5  int `gorm:"not null;default:0"`
	Count6  int `gorm:"not null;default:0"`
	Count7  int `gorm:"not null;default:0"`
	Count8  int `gorm:"not null;default:0"`
	Count9  int `gorm:"not null;default:0"`
	Count10 int `gorm:"not null;default:0"`
	Count11 int `gorm:"not null;default:0"`
	Count12 int `gorm:"not null;default:0"`
	Count13 int `gorm:"not null;default:0"`
	Count14 int `gorm:"not null;default:0"`
	Count15 int `gorm:"not null;default:0"`
	Count16 int `gorm:"not null;default:0"`
	Count17 int `gorm:"not null;default:0"`
	Count18 int `gorm:"not null;default:0"`
	Count19 int `gorm:"not null;default:0"`
	Count20 int `gorm:"not null;default:0"`
}

func (u *Score) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

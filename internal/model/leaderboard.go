package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Leaderboard struct {
	UserID    uuid.UUID    `gorm:"type:uuid;not null;unique;primaryKey"`
	IsTop     pq.BoolArray `gorm:"type:bool[];not null;default:'{false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false}'"`
	UpdatedAt time.Time
}

func (l *Leaderboard) BeforeCreate(tx *gorm.DB) (err error) {
	if l.IsTop != nil && len(l.IsTop) != 20 {
		err = fmt.Errorf("len(Leaderboard.IsTop) must be 20 (UserID = %s)", l.UserID.String())
	}

	return
}

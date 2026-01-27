package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Score struct {
	// auto fk because field name convention + we declare Score field in User
	UserID uuid.UUID `gorm:"type:uuid;not null;unique;primaryKey"`
	User   User
	Count  pq.Int32Array `gorm:"type:int[];not null;default:array_fill(0, ARRAY[20])"`
}

func (u *Score) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Count != nil && len(u.Count) != 20 {
		err = fmt.Errorf("len(Score.Count) must be 20 (UserID = %s)", u.UserID.String())
	}

	return
}

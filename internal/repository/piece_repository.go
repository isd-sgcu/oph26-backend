package repository

import (
	"errors"
	"oph26-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PieceRepository interface {
	FindAttendeeByUserID(userID uuid.UUID) (*entity.Attendee, error)
	FindMyPieceByAttendeeID(attendeeID uuid.UUID) (*entity.MyPiece, error)
}

type PieceRepositoryImpl struct {
	DB *gorm.DB
}

func NewPieceRepository(db *gorm.DB) PieceRepository {
	return &PieceRepositoryImpl{DB: db}
}

func (r *PieceRepositoryImpl) FindAttendeeByUserID(userID uuid.UUID) (*entity.Attendee, error) {
	var attendee entity.Attendee
	if err := r.DB.Where("user_id = ?", userID).First(&attendee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attendee, nil
}

func (r *PieceRepositoryImpl) FindMyPieceByAttendeeID(attendeeID uuid.UUID) (*entity.MyPiece, error) {
	var piece entity.MyPiece
	if err := r.DB.Where("attendee_id = ?", attendeeID).First(&piece).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &piece, nil
}

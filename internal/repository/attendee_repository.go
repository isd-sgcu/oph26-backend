package repository

import (
	"errors"
	"oph26-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendeeRepository interface {
	FindByUserID(userID uuid.UUID) (*entity.Attendee, error)
	FindByTicketCode(ticketCode string) (*entity.Attendee, error)
}

type AttendeeRepositoryImpl struct {
	DB *gorm.DB
}

func NewAttendeeRepository(db *gorm.DB) AttendeeRepository {
	return &AttendeeRepositoryImpl{DB: db}
}

func (r *AttendeeRepositoryImpl) FindByUserID(userID uuid.UUID) (*entity.Attendee, error) {
	var attendee entity.Attendee
	if err := r.DB.Where("user_id = ?", userID).First(&attendee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (r *AttendeeRepositoryImpl) FindByTicketCode(ticketCode string) (*entity.Attendee, error) {
	var attendee entity.Attendee
	err := r.DB.Preload("CheckinStaff").Where("ticket_code = ?", ticketCode).First(&attendee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

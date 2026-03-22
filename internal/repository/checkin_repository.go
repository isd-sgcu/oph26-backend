package repository

import (
	"fmt"
	"oph26-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckinRepository interface {
	FindCheckinByAttendeeAndFaculty(attendeeID uuid.UUID, faculty string) ([]entity.Checkin, error)
	CreateCheckin(attendeeId uuid.UUID, faculty string, staffId uuid.UUID) error
	Checkinstatus(attendeeId uuid.UUID) (bool, error)
}

type CheckinRepositoryImpl struct {
	DB *gorm.DB
}

func NewCheckinRepository(db *gorm.DB) CheckinRepository {
	return &CheckinRepositoryImpl{DB: db}
}

func (r *CheckinRepositoryImpl) FindCheckinByAttendeeAndFaculty(attendeeID uuid.UUID, faculty string) ([]entity.Checkin, error) {
	var checkins []entity.Checkin
	err := r.DB.Where("attendee_id = ? AND faculty = ?", attendeeID, faculty).Find(&checkins).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find check-ins: %w", err)
	}
	return checkins, nil
}

func (r *CheckinRepositoryImpl) CreateCheckin(attendeeId uuid.UUID, faculty string, staffId uuid.UUID) error {
	newCheckin := &entity.Checkin{
		ID:         uuid.New(),
		AttendeeID: attendeeId,
		Faculty:    faculty,
		StaffID:    staffId,
	}
	return r.DB.Create(newCheckin).Error
}

func (r *CheckinRepositoryImpl) Checkinstatus(attendeeId uuid.UUID) (bool, error) {
	var checkin entity.Checkin
	err := r.DB.Where("attendee_id = ?", attendeeId).First(&checkin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check checkin status: %w", err)
	}
	return true, nil
}



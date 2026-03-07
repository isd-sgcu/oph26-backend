package repository

import (
	"oph26-backend/internal/entity"

	"gorm.io/gorm"
)

type AttendeeRepository interface {
	Upsert(attendee *entity.Attendee) (bool, error)
	CountByAttendeeType(attendeeType string) (int64, error)
}

type AttendeeRepositoryImpl struct {
	DB *gorm.DB
}

func NewAttendeeRepository(db *gorm.DB) AttendeeRepository {
	return &AttendeeRepositoryImpl{DB: db}
}

func (r *AttendeeRepositoryImpl) Upsert(attendee *entity.Attendee) (founded bool, err error) {
	result := r.DB.Where("user_id = ?", attendee.UserID).FirstOrCreate(attendee)
	return result.RowsAffected == 0, result.Error
}

func (r *AttendeeRepositoryImpl) CountByAttendeeType(attendeeType string) (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Attendee{}).Where("attendee_type = ?", attendeeType).Count(&count).Error
	return count, err
}

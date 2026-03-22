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
	Update(attendee *entity.Attendee, userId uuid.UUID) error
	Upsert(attendee *entity.Attendee) (bool, error)
	CountByAttendeeType(attendeeType string) (int64, error)
	CreateMyPieceAndLink(attendee *entity.Attendee, myPiece *entity.MyPiece) error
	GetFavWorkshop(userID uuid.UUID) (*entity.StringSet, error)
	UpdateAttendeeRank(userID uuid.UUID) error
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
	if err := r.DB.Where("ticket_code = ?", ticketCode).First(&attendee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (r *AttendeeRepositoryImpl) Upsert(attendee *entity.Attendee) (found bool, err error) {
	result := r.DB.Where("user_id = ?", attendee.UserID).FirstOrCreate(attendee)
	return result.RowsAffected == 0, result.Error
}

func (r *AttendeeRepositoryImpl) CountByAttendeeType(attendeeType string) (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Attendee{}).Where("attendee_type = ?", attendeeType).Count(&count).Error
	return count, err
}

func (r *AttendeeRepositoryImpl) CreateMyPieceAndLink(attendee *entity.Attendee, myPiece *entity.MyPiece) error {
	if err := r.DB.Create(myPiece).Error; err != nil {
		return err
	}
	return nil
}

func (r *AttendeeRepositoryImpl) Update(attendee *entity.Attendee, userId uuid.UUID) error {
	res := r.DB.Model(&entity.Attendee{}).
		Where(&entity.Attendee{UserID: userId}).
		Updates(attendee)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *AttendeeRepositoryImpl) GetFavWorkshop(userID uuid.UUID) (*entity.StringSet, error) {
	var set entity.StringSet
	res := r.DB.Model(&entity.Attendee{}).
		Select("favorite_workshops").
		Where(&entity.Attendee{UserID: userID}).
		Scan(&set)

	if res.Error != nil {
		return nil, res.Error
	}

	return &set, nil
}

func (r *AttendeeRepositoryImpl) UpdateAttendeeRank(userID uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(42)").Error; err != nil {
			return err
		}
		var maxRank int
		if err := tx.Model(&entity.Attendee{}).Select("COALESCE(MAX(rank), 0)").Where("rank > 0").Scan(&maxRank).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Attendee{}).
			Where("user_id = ? AND rank <= 0", userID).
			Update("rank", maxRank+1).Error
	})
}

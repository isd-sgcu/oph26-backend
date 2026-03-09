package repository

import (
	"errors"
	"oph26-backend/internal/entity"

	"gorm.io/gorm"
)

type StaffRepository interface {
	FindByUserID(userID *string) (*entity.Staff, error)
	FindByEmail(email string) (*entity.Staff, error)
	Create(staff *entity.Staff) error
	Update(staff *entity.Staff) error
}

type StaffRepositoryImpl struct {
	DB *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &StaffRepositoryImpl{DB: db}
}

func (r *StaffRepositoryImpl) FindByUserID(userID *string) (*entity.Staff, error) {
	var staff entity.Staff
	if err := r.DB.Where("user_id = ?", userID).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &staff, nil
}

func (r *StaffRepositoryImpl) FindByEmail(email string) (*entity.Staff, error) {
	var staff entity.Staff
	if err := r.DB.Where("email = ?", email).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &staff, nil
}

func (r *StaffRepositoryImpl) Create(staff *entity.Staff) error {
	return r.DB.Create(staff).Error
}

func (r *StaffRepositoryImpl) Update(staff *entity.Staff) error {
	return r.DB.Save(staff).Error
}

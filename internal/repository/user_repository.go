package repository

import (
	"errors"
	"oph26-backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryImpl) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

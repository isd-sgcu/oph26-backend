package repository

import (
	"errors"
	"oph26-backend/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Store(refreshToken *entity.RefreshToken) error
	FindByToken(token string) (*entity.RefreshToken, error)
	Revoke(token string) error
	RevokeAllByUserID(userID uuid.UUID) error
	IsActive(token string) (bool, error)
	DeleteExpired(now time.Time) error
}

type RefreshTokenRepositoryImpl struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{DB: db}
}

func (r *RefreshTokenRepositoryImpl) Store(refreshToken *entity.RefreshToken) error {
	return r.DB.Create(refreshToken).Error
}

func (r *RefreshTokenRepositoryImpl) FindByToken(token string) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	err := r.DB.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) Revoke(token string) error {
	return r.DB.Model(&entity.RefreshToken{}).
		Where("token = ? AND is_revoked = ?", token, false).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepositoryImpl) RevokeAllByUserID(userID uuid.UUID) error {
	return r.DB.Model(&entity.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepositoryImpl) IsActive(token string) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.RefreshToken{}).
		Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RefreshTokenRepositoryImpl) DeleteExpired(now time.Time) error {
	return r.DB.Where("expires_at <= ?", now).Delete(&entity.RefreshToken{}).Error
}

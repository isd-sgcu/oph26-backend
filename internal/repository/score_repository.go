package repository

import (
	"fmt"
	"oph26-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScoreRepositoryImpl struct {
	DB *gorm.DB
}

type ScoreRepository interface {
	IncrementCountByIndex(userID uuid.UUID, index int) error
	FindAll() ([]entity.Score, error)
	Count() (int, error)
	Create(score *entity.Score) error
}

func NewScoreRepository(db *gorm.DB) ScoreRepository {
	return &ScoreRepositoryImpl{DB: db}
}

func (r *ScoreRepositoryImpl) IncrementCountByIndex(userID uuid.UUID, index int) error {
	if index < 1 || index > 20 {
		return fmt.Errorf("faculty index must be in range 1..20, got %d", index)
	}

	column := fmt.Sprintf("count%d", index)
	return r.DB.Model(&entity.Score{}).Where("user_id = ?", userID).Update(column, gorm.Expr(column+" + 1")).Error
}

func (r *ScoreRepositoryImpl) Count() (int, error) {
	var count int64
	err := r.DB.Model(&entity.Score{}).Count(&count).Error

	return int(count), err
}

func (r *ScoreRepositoryImpl) FindAll() ([]entity.Score, error) {
	var scores []entity.Score
	err := r.DB.Find(&scores).Error
	if err != nil {
		return nil, err
	}

	return scores, nil
}

func (r *ScoreRepositoryImpl) Create(score *entity.Score) error {
	return r.DB.Create(score).Error
}

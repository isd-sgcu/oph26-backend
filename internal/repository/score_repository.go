package repository

import (
	"fmt"
	"oph26-backend/internal/entity"
	"reflect"

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
	GetMissingCounts(userID uuid.UUID) (map[int]int, error)
}

func NewScoreRepository(db *gorm.DB) ScoreRepository {
	return &ScoreRepositoryImpl{DB: db}

}

func (r *ScoreRepositoryImpl) IncrementCountByIndex(userID uuid.UUID, index int) error {
	if index < 1 || index > 20 {
		return fmt.Errorf("faculty index must be in range 1..20, got %d", index)
	}

	column := fmt.Sprintf("count%d", index)
	tx := r.DB.Model(&entity.Score{}).Where("user_id = ?", userID).Update(column, gorm.Expr(column+" + 1"))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
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

func (r *ScoreRepositoryImpl) GetMissingCounts(userID uuid.UUID) (map[int]int, error) {
	var score entity.Score
	if err := r.DB.Where("user_id = ?", userID).First(&score).Error; err != nil {
		return nil, err
	}

	result := make(map[int]int)
	scoreValue := reflect.ValueOf(&score).Elem()
	for i := 1; i <= 20; i++ {
		fieldName := fmt.Sprintf("Count%d", i)
		field := scoreValue.FieldByName(fieldName)
		if field.IsValid() && field.Kind() == reflect.Int && field.Int() == 0 {
			var count int64
			column := fmt.Sprintf("count%d", i)
			err := r.DB.Model(&entity.Score{}).Where(column+" = 0 AND user_id != ?", userID).Count(&count).Error
			if err != nil {
				return nil, err
			}
			result[i] = int(count)
		}
	}
	return result, nil
}

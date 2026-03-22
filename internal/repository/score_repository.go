package repository

import (
	"fmt"
	"oph26-backend/internal/entity"
	"reflect"
	"strings"

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
	GetMissingCounts(userID uuid.UUID) (map[int]float64, error)
	IsComplete(userID uuid.UUID) (bool, error)
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

func (r *ScoreRepositoryImpl) GetMissingCounts(userID uuid.UUID) (map[int]float64, error) {
	total, err := r.Count()
	if err != nil {
		return nil, err
	}
	if total <= 1 {
		return make(map[int]float64), nil
	}

	var score entity.Score
	if err := r.DB.Where("user_id = ?", userID).First(&score).Error; err != nil {
		return nil, err
	}

	zeroColumns := make([]int, 0)
	scoreValue := reflect.ValueOf(&score).Elem()
	for i := 1; i <= 20; i++ {
		fieldName := fmt.Sprintf("Count%d", i)
		field := scoreValue.FieldByName(fieldName)
		if field.IsValid() && field.Kind() == reflect.Int && field.Int() == 0 {
			zeroColumns = append(zeroColumns, i)
		}
	}

	if len(zeroColumns) == 0 {
		return make(map[int]float64), nil
	}

	selectParts := make([]string, len(zeroColumns))
	for i, col := range zeroColumns {
		colName := fmt.Sprintf("count%d", col)
		selectParts[i] = fmt.Sprintf("COALESCE(SUM(CASE WHEN %s = 0 THEN 1 ELSE 0 END), 0) AS zero_%d", colName, col)
	}

	type zeroResult struct {
		Column int
		Count  int64
	}

	query := fmt.Sprintf("SELECT %s FROM scores WHERE user_id != ?", strings.Join(selectParts, ", "))
	rows, err := r.DB.Raw(query, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]float64)
	if rows.Next() {
		vals := make([]interface{}, len(zeroColumns))
		ptrs := make([]int64, len(zeroColumns))
		for i := range vals {
			vals[i] = &ptrs[i]
		}
		if err := rows.Scan(vals...); err != nil {
			return nil, err
		}
		for i, col := range zeroColumns {
			percent := (float64(ptrs[i]) / float64(total-1)) * 100
			result[col] = percent
		}
	}

	return result, nil
}

func (r *ScoreRepositoryImpl) IsComplete(userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Score{}).
		Where("user_id = ? AND count1 > 0 AND count2 > 0 AND count3 > 0 AND count4 > 0 AND count5 > 0 AND count6 > 0 AND count7 > 0 AND count8 > 0 AND count9 > 0 AND count10 > 0 AND count11 > 0 AND count12 > 0 AND count13 > 0 AND count14 > 0 AND count15 > 0 AND count16 > 0 AND count17 > 0 AND count18 > 0 AND count19 > 0 AND count20 > 0", userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

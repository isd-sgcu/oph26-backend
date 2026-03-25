package repository

import (
	"fmt"
	"oph26-backend/internal/entity"

	"gorm.io/gorm"
)

type QuestionnaireRepository interface {
	Create(q *entity.Questionnaire) (*entity.Questionnaire, error)
	Update(q *entity.Questionnaire) (*entity.Questionnaire, error)
	ExistsByUserID(userID string) (bool, error)
}

type QuestionnaireRepositoryImpl struct {
	DB *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) QuestionnaireRepository {
	return &QuestionnaireRepositoryImpl{DB: db}
}

func (r *QuestionnaireRepositoryImpl) Create(q *entity.Questionnaire) (*entity.Questionnaire, error) {
	err := r.DB.Create(q).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create questionnaire: %w", err)
	}
	return q, nil
}

func (r *QuestionnaireRepositoryImpl) Update(q *entity.Questionnaire) (*entity.Questionnaire, error) {
	err := r.DB.Save(q).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update questionnaire: %w", err)
	}
	return q, nil
}

func (r *QuestionnaireRepositoryImpl) ExistsByUserID(userID string) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Questionnaire{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check questionnaire existence: %w", err)
	}
	return count > 0, nil
}

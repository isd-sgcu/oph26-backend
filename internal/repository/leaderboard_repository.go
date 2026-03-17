package repository

import (
	"errors"
	"oph26-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaderboardRepositoryImpl struct {
	DB *gorm.DB
}

type LeaderboardRepository interface {
	FindLeaderboardByUserID(userID uuid.UUID) (*entity.Leaderboard, error)
	UpdateIsTop() error
}

func NewLeaderboardRepository(db *gorm.DB) LeaderboardRepository {
	return &LeaderboardRepositoryImpl{DB: db}
}

func (r *LeaderboardRepositoryImpl) FindLeaderboardByUserID(userID uuid.UUID) (*entity.Leaderboard, error) {
	var leaderboard entity.Leaderboard
	if err := r.DB.Where("user_id = ?", userID).First(&leaderboard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &leaderboard, nil
}

func (r *LeaderboardRepositoryImpl) UpdateIsTop() error {
	query := `
	WITH thresholds AS (
		SELECT array_agg(threshold ORDER BY i) AS t
		FROM (
			SELECT
				i,
				percentile_disc(0.99) WITHIN GROUP (ORDER BY s."count"[i]) AS threshold
			FROM scores s
			CROSS JOIN generate_series(1,20) AS i
			GROUP BY i
		) sub
	),
	marks AS (
		SELECT
			user_id,
			ARRAY(
				SELECT s."count"[gs] >= t.t[gs]
				FROM generate_series(1,20) AS gs
			) AS is_top
		FROM scores s, thresholds t
		WHERE EXISTS (
			SELECT 1 FROM generate_series(1,20) gs
			WHERE s."count"[gs] >= t.t[gs]
		)
	)
	UPDATE leaderboards l
	SET is_top = m.is_top
	FROM marks m
	WHERE l.user_id = m.user_id;
	`
	return r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.Leaderboard{}).Update("is_top", gorm.Expr("array_fill(false, ARRAY[20])")).Error
		if err != nil {
			return err
		}

		err = tx.Exec(query).Error
		if err != nil {
			return err
		}

		return nil
	})
}

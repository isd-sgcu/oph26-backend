package usecase

import (
	"oph26-backend/internal/model"
	"oph26-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeaderboardUsecaseImpl struct {
	leaderboardRepo repository.LeaderboardRepository
	scoreRepo       repository.ScoreRepository
}

type LeaderboardUsecase interface {
	GetMyLeaderboard(c *fiber.Ctx) error
	UpdateScore(userID uuid.UUID, facultyIndex int) error
	UpdateLeaderboard() error
}

func NewLeaderboardUsecase(leaderboardRepository repository.LeaderboardRepository, scoreRepository repository.ScoreRepository) LeaderboardUsecase {
	return &LeaderboardUsecaseImpl{leaderboardRepo: leaderboardRepository, scoreRepo: scoreRepository}
}

func (u *LeaderboardUsecaseImpl) GetMyLeaderboard(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user id from context",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	leaderboard, err := u.leaderboardRepo.FindByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if leaderboard == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "leaderboard not found",
		})
	}

	return c.JSON(&model.LeaderboardResponse{
		IsTop: leaderboard.IsTop,
	})
}

func (u *LeaderboardUsecaseImpl) UpdateScore(userID uuid.UUID, facultyIndex int) error {
	return u.scoreRepo.IncrementCountByIndex(userID, facultyIndex)
}

func (u *LeaderboardUsecaseImpl) UpdateLeaderboard() error {
	return u.leaderboardRepo.UpdateIsTop()
}
package usecase

import (
	"oph26-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type StatsUsecase interface {
	GetCountAttendeesStats(c *fiber.Ctx) error
	GetCountUniqueAttendeesCheckinsGroupedByDateStats(c *fiber.Ctx) error
	GetCountAvailablePiecesGroupedByFacultyStats(c *fiber.Ctx) error
}

type statsUsecaseImpl struct {
	statsRepo repository.StatsRepository
}

func NewStatsUsecase(statsRepo repository.StatsRepository) StatsUsecase {
	return &statsUsecaseImpl{
		statsRepo: statsRepo,
	}
}

func (u *statsUsecaseImpl) GetCountAttendeesStats(c *fiber.Ctx) error {
	count, err := u.statsRepo.CountAttendees()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}

func (u *statsUsecaseImpl) GetCountUniqueAttendeesCheckinsGroupedByDateStats(c *fiber.Ctx) error {
	counts, err := u.statsRepo.CountUniqueAttendeesCheckinsGroupedByDate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"counts": counts,
	})
}

func (u *statsUsecaseImpl) GetCountAvailablePiecesGroupedByFacultyStats(c *fiber.Ctx) error {
	countByFaculty, err := u.statsRepo.CountAvailablePiecesGroupedByFaculty()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count_by_faculty": countByFaculty,
	})
}

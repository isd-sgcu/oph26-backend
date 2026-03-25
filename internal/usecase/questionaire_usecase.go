// filepath: c:\Users\Thukdanai Thaothawin\myapp\oph26-backend\internal\usecase\questionaire_usecase.go
package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"oph26-backend/internal/entity"
	"oph26-backend/internal/repository"
)

type QuestionnaireUsecase interface {
	CreateQuestionnaire(c *fiber.Ctx) error
	GetQuestionnaire(c *fiber.Ctx) error
}

type questionnaireUsecase struct {
	repo repository.QuestionnaireRepository
}

func NewQuestionnaireUsecase(repo repository.QuestionnaireRepository) QuestionnaireUsecase {
	return &questionnaireUsecase{repo: repo}
}

func (u *questionnaireUsecase) CreateQuestionnaire(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user id from context",
		})
	}

	var answers datatypes.JSON
	if err := c.BodyParser(&answers); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	q := &entity.Questionnaire{
		ID:      userID,
		UserID:  userID.String(),
		Answers: answers,
	}

	if exists, err := u.repo.ExistsByUserID(userID.String()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	} else if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Questionnaire already exists for this user",
		})
	}

	questionnaire, err := u.repo.Create(q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(questionnaire)
}

func (u *questionnaireUsecase) GetQuestionnaire(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user id from context",
		})
	}

	exists, err := u.repo.ExistsByUserID(userID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exists": exists,
	})
}

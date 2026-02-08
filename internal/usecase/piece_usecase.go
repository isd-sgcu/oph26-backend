package usecase

import (
	pieceModel "oph26-backend/internal/model/piece"
	"oph26-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PieceUsecase interface {
	GetMyPiece(c *fiber.Ctx) error
}

type PieceUsecaseImpl struct {
	PieceRepo repository.PieceRepository
}

func NewPieceUsecase(pieceRepo repository.PieceRepository) PieceUsecase {
	return &PieceUsecaseImpl{
		PieceRepo: pieceRepo,
	}
}

func (u *PieceUsecaseImpl) GetMyPiece(c *fiber.Ctx) error {
	// TODO: Auth
	role, _ := c.Locals("role").(string)
	if role == "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, staff accounts cannot access pieces",
		})
	}

	userIDStr, _ := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user ID",
		})
	}

	attendee, err := u.PieceRepo.FindAttendeeByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pieces not found for the current user",
		})
	}

	piece, err := u.PieceRepo.FindMyPieceByAttendeeID(attendee.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if piece == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pieces not found for the current user",
		})
	}

	return c.JSON(pieceModel.MyPieceResponse{
		ID:         piece.ID,
		UserID:     attendee.UserID,
		PieceCode:  piece.PieceCode,
		ExpireDate: piece.ExpireDate,
		Faculty:    attendee.InitialFirstInterestedFaculty,
	})
}

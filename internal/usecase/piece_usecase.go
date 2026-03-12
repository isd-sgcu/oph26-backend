package usecase

import (
	"errors"
	pieceModel "oph26-backend/internal/model/piece"
	"oph26-backend/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PieceUsecase interface {
	GetMyPiece(c *fiber.Ctx) error
	GetCollectedPieces(c *fiber.Ctx) error
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

	if time.Now().After(piece.ExpireDate) {
		newCode, err := generatePieceCode()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot generate new piece code",
			})
		}
		if err := u.PieceRepo.RefreshMyPiece(piece, newCode); err != nil {
		const maxRetries = 5
		var newPiece *entity.MyPiece
		for range maxRetries {
			newCode, err := generatePieceCode()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Cannot generate new piece code",
				})
			}
			candidate := &entity.MyPiece{
				AttendeeID: piece.AttendeeID,
				PieceCode:  newCode,
				ExpireDate: time.Now().Add(24 * time.Hour),
			}
			err = u.PieceRepo.CreateMyPiece(candidate)
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				continue
			}
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create new piece",
				})
			}
			newPiece = candidate
			break
		}
		if newPiece == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to refresh piece",
				"error": "Failed to create new piece after retries",
			})
		}
	}

	return c.JSON(pieceModel.MyPieceResponse{
		ID:         piece.ID,
		UserID:     attendee.UserID,
		PieceCode:  piece.PieceCode,
		ExpireDate: piece.ExpireDate,
		Faculty:    attendee.InitialFirstInterestedFaculty,
	})
}

func (u *PieceUsecaseImpl) GetCollectedPieces(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role == "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, staff accounts cannot access collected pieces",
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
			"error": "No collected pieces found for the current user",
		})
	}

	collected, err := u.PieceRepo.FindCollectedPiecesByAttendeeID(attendee.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	friendPieces := make([]pieceModel.FriendPieceResponse, 0, len(collected))
	for _, cp := range collected {
		fp := pieceModel.FriendPieceResponse{
			ID:          cp.PieceID,
			UserID:      cp.MyPiece.Attendee.UserID,
			Faculty:     cp.MyPiece.Attendee.InitialFirstInterestedFaculty,
			CollectedAt: &cp.CollectedAt,
		}
		friendPieces = append(friendPieces, fp)
	}

	facultyCounts, err := u.PieceRepo.CountCollectedByFaculty(attendee.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	thresholds, err := u.PieceRepo.CountTop1ThresholdByFaculty()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	collectedByFaculty := make(map[string]pieceModel.FacultyStats)
	for faculty, count := range facultyCounts {
		threshold, ok := thresholds[faculty]
		isTop1 := ok && threshold > 0 && count >= threshold
		collectedByFaculty[faculty] = pieceModel.FacultyStats{
			Count:  count,
			IsTop1: isTop1,
		}
	}

	totalCollected := len(collected)

	return c.JSON(pieceModel.CollectedPiecesResponse{
		CollectedPieces: friendPieces,
		Stats: pieceModel.CollectedPiecesStats{
			TotalCollected:     totalCollected,
			CollectedByFaculty: collectedByFaculty,
		},
	})
}

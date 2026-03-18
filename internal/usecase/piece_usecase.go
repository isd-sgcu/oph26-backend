package usecase

import (
	"errors"
	pieceModel "oph26-backend/internal/model/piece"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/repository"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var facultyIndex = map[string]int{
	"edu": 1, "psy": 2, "dent": 3, "law": 4, "commarts": 5,
	"cbs": 6, "md": 7, "pharm": 8, "polsci": 9, "sci": 10,
	"spsc": 11, "eng": 12, "faa": 13, "econ": 14, "arch": 15,
	"ahs": 16, "vet": 17, "arts": 18, "scii": 19, "cusar": 20,
}

type PieceUsecase interface {
	GetMyPiece(c *fiber.Ctx) error
	GetCollectedPieces(c *fiber.Ctx) error
	CollectPiece(c *fiber.Ctx) error
}

type PieceUsecaseImpl struct {
	PieceRepo       repository.PieceRepository
	LeaderboardCase LeaderboardUsecase
	validate        *validator.Validate
}

func NewPieceUsecase(pieceRepo repository.PieceRepository, leaderboardUsecase LeaderboardUsecase) PieceUsecase {
	return &PieceUsecaseImpl{
		PieceRepo:       pieceRepo,
		LeaderboardCase: leaderboardUsecase,
		validate:        validator.New(),
	}
}

func (u *PieceUsecaseImpl) GetMyPiece(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role == "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, staff accounts cannot access pieces",
		})
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user ID",
		})
	}

	attendee, err := u.PieceRepo.FindAttendeeByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pieces not found for the current user",
		})
	}
	if attendee.AttendeeType != "student" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, only student attendees can access pieces",
		})
	}

	piece, err := u.PieceRepo.FindMyPieceByAttendeeID(attendee.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user ID",
		})
	}

	attendee, err := u.PieceRepo.FindAttendeeByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No collected pieces found for the current user",
		})
	}

	collected, err := u.PieceRepo.FindCollectedPiecesByAttendeeID(attendee.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	thresholds, err := u.PieceRepo.CountTop1ThresholdByFaculty()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	return c.JSON(pieceModel.CollectedPiecesResponse{
		CollectedPieces: friendPieces,
		Stats: pieceModel.CollectedPiecesStats{
			TotalCollected:     len(collected),
			CollectedByFaculty: collectedByFaculty,
		},
	})
}

func (u *PieceUsecaseImpl) CollectPiece(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role == "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, staff accounts cannot collect pieces",
		})
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid user ID",
		})
	}

	attendee, err := u.PieceRepo.FindAttendeeByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendee not found",
		})
	}

	var reqBody pieceModel.CollectPieceRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := u.validate.Struct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; piece_code is required",
		})
	}

	matched, _ := regexp.MatchString(`^[A-Z0-9]{6}$`, reqBody.PieceCode)
	if !matched {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid piece code",
		})
	}

	friendPiece, err := u.PieceRepo.FindMyPieceByCode(reqBody.PieceCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if friendPiece == nil || friendPiece.ExpireDate.Before(time.Now()) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Piece not found or expired",
		})
	}

	if friendPiece.AttendeeID == attendee.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot collect your own piece",
		})
	}

	existing, err := u.PieceRepo.FindCollectedPiece(attendee.ID, friendPiece.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if existing != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Already collected this piece",
		})
	}

	now := time.Now()
	cp := entity.CollectedPiece{
		AttendeeID:  attendee.ID,
		PieceID:     friendPiece.ID,
		CollectedAt: now,
	}
	if err := u.PieceRepo.CreateCollectedPiece(&cp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	faculty := friendPiece.Attendee.InitialFirstInterestedFaculty
	if idx, ok := facultyIndex[string(faculty)]; ok {
		_ = u.LeaderboardCase.UpdateScore(userID, idx)
		_ = u.LeaderboardCase.UpdateLeaderboard()
	}

	return c.JSON(pieceModel.CollectPieceResponse{
		Ok: true,
		CollectedPiece: pieceModel.FriendPieceResponse{
			ID:          friendPiece.ID,
			UserID:      friendPiece.Attendee.UserID,
			Faculty:     faculty,
			CollectedAt: &now,
		},
	})
}

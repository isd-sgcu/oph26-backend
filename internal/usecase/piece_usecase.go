package usecase

import (
	"oph26-backend/internal/entity"
	pieceModel "oph26-backend/internal/model/piece"
	"oph26-backend/internal/repository"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var facultyIndex = map[string]int{
	"edu": 1, "psy": 2, "dent": 3, "law": 4, "commarts": 5,
	"cbs": 6, "md": 7, "pharm": 8, "polsci": 9, "sci": 10,
	"spsc": 11, "eng": 12, "faa": 13, "econ": 14, "arch": 15,
	"ahs": 16, "vet": 17, "arts": 18, "scii": 19, "cusar": 20,
}

var pieceCodeRegex = regexp.MustCompile(`^[A-Z0-9]{6}$`)

type PieceUsecase interface {
	GetMyPiece(c *fiber.Ctx) error
	GetCollectedPieces(c *fiber.Ctx) error
	CollectPiece(c *fiber.Ctx) error
}

type PieceUsecaseImpl struct {
	AttendeeRepo    repository.AttendeeRepository
	PieceRepo       repository.PieceRepository
	LeaderboardCase LeaderboardUsecase
	ScoreRepo       repository.ScoreRepository
	validate        *validator.Validate
}

func NewPieceUsecase(pieceRepo repository.PieceRepository, leaderboardUsecase LeaderboardUsecase, scoreRepo repository.ScoreRepository, attendeeRepo repository.AttendeeRepository) PieceUsecase {
	return &PieceUsecaseImpl{
		AttendeeRepo:    attendeeRepo,
		PieceRepo:       pieceRepo,
		LeaderboardCase: leaderboardUsecase,
		ScoreRepo:       scoreRepo,
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

	// If the piece is expired, update the code and expiration date
	if piece.ExpireDate.Before(time.Now()) {
		maxRetries := 5
		for range maxRetries {
			newCode, pErr := generatePieceCode()
			if pErr != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": pErr.Error()})
			}
			piece, err := u.PieceRepo.RefreshMyPiece(piece, newCode)
			if err == nil {
				return c.JSON(pieceModel.MyPieceResponse{
					ID:         piece.ID,
					UserID:     attendee.UserID,
					PieceCode:  piece.PieceCode,
					ExpireDate: piece.ExpireDate,
					Faculty:    attendee.InitialFirstInterestedFaculty,
				})
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh piece after multiple attempts, please try again later",
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

	sameMissingCount, err := u.ScoreRepo.GetMissingCounts(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pieceModel.CollectedPiecesResponse{
		CollectedPieces: friendPieces,
		Stats: pieceModel.CollectedPiecesStats{
			TotalCollected:     len(collected),
			CollectedByFaculty: collectedByFaculty,
			SameMissingCount:   sameMissingCount,
			Rank:               attendee.Rank,
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

	if !pieceCodeRegex.MatchString(reqBody.PieceCode) {
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
	idx, ok := facultyIndex[string(faculty)]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unknown faculty",
		})
	}

	if err := u.ScoreRepo.IncrementCountByIndex(userID, idx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	
	if attendee.Rank <= 0 {
		isComplete, err := u.ScoreRepo.IsComplete(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if isComplete {
			if err := u.AttendeeRepo.UpdateAttendeeRank(attendee.UserID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
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

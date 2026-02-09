package usecase

import (
	"crypto/rand"
	"math/big"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model/user"
	"oph26-backend/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AttendeeUsecase interface {
	PostAttendeesUseCase(c *fiber.Ctx) error
}

type AttendeeUsecaseImpl struct {
	UserRepository     repository.UserRepository
	AttendeeRepository repository.AttendeeRepository
}

func NewAttendeeUsecase(userRepository repository.UserRepository, attendeeRepository repository.AttendeeRepository) AttendeeUsecase {
	return &AttendeeUsecaseImpl{
		UserRepository:     userRepository,
		AttendeeRepository: attendeeRepository,
	}
}

func (u *AttendeeUsecaseImpl) PostAttendeesUseCase(c *fiber.Ctx) error {
	request := new(user.AttendeeCreateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// this wont ever be invalid right?
	userIdRaw := c.Locals("user_id").(string)
	userId, err := uuid.Parse(userIdRaw)
	// userId := uuid.New()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id",
		})
	}

	// staff only
	role := c.Locals("role").(string)
	if role == "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "This is for self-registration. (non staff only)",
		})
	}

	code, _ := generatePieceCode()

	attendee := entity.Attendee{
		UserID:                        userId,
		Firstname:                     request.Firstname,
		Surname:                       request.Surname,
		AttendeeType:                  request.AttendeeType,
		Age:                           request.Age,
		Province:                      request.Province,
		StudyLevel:                    request.StudyLevel,
		SchoolName:                    request.SchoolName,
		NewsSourceSelected:            request.NewsSourceSelected,
		NewsSourcesOther:              request.NewsSourcesOther,
		InterestedFaculty:             request.InterestedFaculty,
		InitialFirstInterestedFaculty: request.InterestedFaculty[0],
		ObjectiveSelected:             request.ObjectiveSelected,
		ObjectiveOther:                request.ObjectiveOther,
	}

	if request.AttendeeType == "highschool" {
		attendee.MyPiece = &entity.MyPiece{
			PieceCode:  code,
			ExpireDate: time.Now().Add(24 * time.Hour),
		}
	}

	founded, err2 := u.AttendeeRepository.Upsert(&attendee)
	if err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal DB error",
		})
	}
	if founded {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Attendee already exists",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok": true,
	})
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ^[A-Z0-9]{6}$
func generatePieceCode() (string, error) {
	b := make([]byte, 6)
	for i := range b {
		index, err := rand.Int(rand.Reader, big.NewInt(6))
		if err != nil {
			return "", err
		}
		b[i] = charset[index.Int64()]
	}

	return string(b), nil
}

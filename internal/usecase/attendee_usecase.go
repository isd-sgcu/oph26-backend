package usecase

import (
	"crypto/rand"
	"math/big"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model"
	"oph26-backend/internal/model/attendee"
	"oph26-backend/internal/repository"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
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
	// userId := uuid.New()

	request := new(attendee.AttendeeCreateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validationnnnnn
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Fail to validate JSON",
		})
	}

	if !model.NewsSourcesAreValid(request.NewsSourceSelected) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknowns new source",
		})
	}

	if slices.Contains(request.NewsSourceSelected, string(model.OtherNewsSource)) && request.NewsSourcesOther == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; news_sources_selected is 'อื่น ๆ', but news_sources_other is not provided",
		})
	}

	if !model.ObjectivesAreValid(request.ObjectiveSelected) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown objective",
		})
	}

	if slices.Contains(request.ObjectiveSelected, string(model.OtherObjective)) && request.ObjectiveOther == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; objective_selected is 'อื่น ๆ', but objective_other is not provided",
		})
	}

	if !model.ProvinceIsValid(request.Province) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown province",
		})
	}

	if request.StudyLevel != nil && !model.StudyLevelIsValid(*request.StudyLevel) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown study_level",
		})
	}

	// validate options array
	arr := []string(request.InterestedFaculty)
	if !model.FacultiesAreValid(arr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown faculty",
		})
	}

	arr = []string(request.ObjectiveSelected)
	if !model.ObjectivesAreValid(arr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown objective",
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
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[index.Int64()]
	}

	return string(b), nil
}

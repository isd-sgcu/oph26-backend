package usecase

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model"
	attendeeModel "oph26-backend/internal/model/attendee"
	"oph26-backend/internal/repository"
	"regexp"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendeesUsecase interface {
	GetMyAttendee(c *fiber.Ctx) error
	GetByAttendeeId(c *fiber.Ctx) error
	PostAttendee(c *fiber.Ctx) error
}

type AttendeesUsecaseImpl struct {
	AttendeeRepository repository.AttendeeRepository
	UserRepository     repository.UserRepository
}

func NewAttendeeUsecase(attendeeRepositry repository.AttendeeRepository, userRepository repository.UserRepository) AttendeesUsecase {
	return &AttendeesUsecaseImpl{
		AttendeeRepository: attendeeRepositry,
		UserRepository:     userRepository,
	}
}

func (u *AttendeesUsecaseImpl) GetMyAttendee(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user role from context",
		})
	}
	// TODO: Auth here
	if role == "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, staff accounts cannot access attendee data",
		})
	}

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

	attendee, err := u.AttendeeRepository.FindByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendee data not found for the current user",
		})
	}

	return c.JSON(&attendeeModel.AttendeeResponse{
		Age:                           attendee.Age,
		AttendeeType:                  attendee.AttendeeType,
		CertificateName:               attendee.CertificateName,
		CheckedInAt:                   attendee.CheckedInAt,
		CheckinStaffID:                attendee.CheckinStaffID,
		CreatedAt:                     attendee.CreatedAt,
		FavoriteWorkshops:             attendee.FavoriteWorkshops,
		Firstname:                     attendee.Firstname,
		ID:                            attendee.ID,
		InitialFirstInterestedFaculty: attendee.InitialFirstInterestedFaculty,
		InterestedFaculty:             attendee.InterestedFaculty,
		NewsSourcesOther:              attendee.NewsSourcesOther,
		NewsSourceSelected:            attendee.NewsSourceSelected,
		ObjectiveOther:                attendee.ObjectiveOther,
		ObjectiveSelected:             attendee.ObjectiveSelected,
		Province:                      attendee.Province,
		SchoolName:                    attendee.SchoolName,
		StudyLevel:                    attendee.StudyLevel,
		Surname:                       attendee.Surname,
		TicketCode:                    attendee.TicketCode,
		UpdatedAt:                     attendee.UpdatedAt,
		UserID:                        attendee.UserID,
	})
}

func (u *AttendeesUsecaseImpl) GetByAttendeeId(c *fiber.Ctx) error {
	ticketCode := c.Params("attendeeId")
	matched, _ := regexp.MatchString(`^[HSPEA]\d{6}$`, ticketCode)
	if !matched {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Ticket Code",
		})
	}

	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user role from context",
		})
	}
	// TODO: Auth here
	if role == "attendee" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	attendee, err := u.AttendeeRepository.FindByTicketCode(ticketCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendee not found",
		})
	}

	var checkinStaff *model.StaffResponse
	if attendee.CheckinStaff != nil {
		staff := attendee.CheckinStaff
		checkinStaff = &model.StaffResponse{
			ID:        staff.ID,
			UserID:    staff.UserID,
			Cuid:      staff.Cuid,
			Firstname: staff.Firstname,
			Surname:   staff.Surname,
			Nickname:  staff.Nickname,
			Phone:     staff.Phone,
			Year:      staff.Year,
			Email:     staff.Email,
			Faculty:   staff.Faculty,
			CreatedAt: staff.CreatedAt,
			UpdatedAt: staff.UpdatedAt,
		}
	}

	return c.JSON(&attendeeModel.AttendeeStaffResponse{
		AttendeeResponse: attendeeModel.AttendeeResponse{
			Age:                           attendee.Age,
			AttendeeType:                  attendee.AttendeeType,
			CertificateName:               attendee.CertificateName,
			CheckedInAt:                   attendee.CheckedInAt,
			CheckinStaffID:                attendee.CheckinStaffID,
			CreatedAt:                     attendee.CreatedAt,
			FavoriteWorkshops:             attendee.FavoriteWorkshops,
			Firstname:                     attendee.Firstname,
			ID:                            attendee.ID,
			InitialFirstInterestedFaculty: attendee.InitialFirstInterestedFaculty,
			InterestedFaculty:             attendee.InterestedFaculty,
			NewsSourcesOther:              attendee.NewsSourcesOther,
			NewsSourceSelected:            attendee.NewsSourceSelected,
			ObjectiveOther:                attendee.ObjectiveOther,
			ObjectiveSelected:             attendee.ObjectiveSelected,
			Province:                      attendee.Province,
			SchoolName:                    attendee.SchoolName,
			StudyLevel:                    attendee.StudyLevel,
			Surname:                       attendee.Surname,
			TicketCode:                    attendee.TicketCode,
			UpdatedAt:                     attendee.UpdatedAt,
			UserID:                        attendee.UserID,
		},
		CheckinStaff: checkinStaff,
	})
}

func (u *AttendeesUsecaseImpl) PostAttendee(c *fiber.Ctx) error {
	userIdRaw, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id",
		})
	}

	userId, err := uuid.Parse(userIdRaw)
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

	request := new(attendeeModel.AttendeeCreateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validation
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Fail to validate JSON",
		})
	}

	if !model.NewsSourcesAreValid(request.NewsSourceSelected) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown news source",
		})
	}

	if slices.Contains(request.NewsSourceSelected, string(model.OtherNewsSource)) && request.NewsSourcesOther == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; news_sources_selected is 'อื่น ๆ', but news_sources_other is not provided",
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

	code, codeErr := generatePieceCode()
	if codeErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot generate piece code",
		})
	}

	// Generate ticket code, please refer to the docs (บรีฟเว็บไซต์ section)
	// note that their might be multiple user doing this simultaneously, so
	// we are just going to loop until its successfully created
	retryCount := 5 // exponential backoff???

	for retryCount > 0 {
		retryCount -= 1

		ticketCode, ticketCodeErr := u.generateTicketCode(request.AttendeeType)
		if ticketCodeErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal DB error",
			})
		}

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
			TicketCode:                    ticketCode,
		}

		if request.AttendeeType == "highschool" {
			attendee.MyPiece = &entity.MyPiece{
				PieceCode:  code,
				ExpireDate: time.Now().Add(24 * time.Hour),
			}
		}

		founded, err2 := u.AttendeeRepository.Upsert(&attendee)
		// TODO: this need `TranslateError: true`
		// - also test this
		if errors.Is(err2, gorm.ErrDuplicatedKey) {
			continue
		}
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

func (u *AttendeesUsecaseImpl) generateTicketCode(attendeeType string) (string, error) {
	ticketCodePrefix := "A"
	switch attendeeType {
	case "elementaryschool":
		ticketCodePrefix = "S"
	case "highschool":
		ticketCodePrefix = "H"
	case "parent":
		ticketCodePrefix = "P"
	case "educationstaff":
		ticketCodePrefix = "E"
	case "other":
		ticketCodePrefix = "A"
	}

	count, err := u.AttendeeRepository.CountByAttendeeType(attendeeType)
	if err != nil {
		return "", err
	}

	return ticketCodePrefix + fmt.Sprintf("%06d", count+1), nil
}

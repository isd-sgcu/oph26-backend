package usecase

import (
	checkinModel "oph26-backend/internal/model/checkin"
	"oph26-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type CheckinUsecase interface {
	CheckIn(c *fiber.Ctx) error
}

type CheckinUsecaseImpl struct {
	AttendeeRepository repository.AttendeeRepository
	StaffRepository    repository.StaffRepository
	CheckinRepository  repository.CheckinRepository
}

func NewCheckinUsecase(attendeeRepository repository.AttendeeRepository, staffRepository repository.StaffRepository, checkinRepository repository.CheckinRepository) CheckinUsecase {
	return &CheckinUsecaseImpl{
		AttendeeRepository: attendeeRepository,
		StaffRepository:    staffRepository,
		CheckinRepository:  checkinRepository,
	}
}

func (u *CheckinUsecaseImpl) CheckIn(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role != "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	userId, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req checkinModel.CheckiAttendeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	attendee, attendeeErr := u.AttendeeRepository.FindByTicketCode(req.TicketCode)
	if attendeeErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find attendee",
		})
	}

	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendee not found",
		})
	}

	staff, staffErr := u.StaffRepository.FindByUserID(userId)
	if staffErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find staff information",
		})
	}

	checkins, checkinsFindErr := u.CheckinRepository.FindCheckinByAttendeeAndFaculty(attendee.ID, staff.Faculty)
	if checkinsFindErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check check-in status",
		})
	}

	// There's already a check-in record for this attendee and faculty, so we return a conflict response
	// The response body is compliant with api-spec.yml definition for CheckinConflictResponse
	if len(checkins) > 0 {
		firstCheckin := checkins[0]
		conflictResponseBody := checkinModel.CheckinConflictResponse{
			Error: "Attendee already checked in with this faculty",
			CheckinResponse: checkinModel.CheckinResponse{
				CheckedInAt: firstCheckin.CheckedInAt,
				UserID:      attendee.UserID,
				Firstname:   attendee.Firstname,
				Surname:     attendee.Surname,
				TicketCode:  attendee.TicketCode,
				Faculty:     staff.Faculty,
			},
		}
		return c.Status(fiber.StatusConflict).JSON(conflictResponseBody)
	}

	// Actually create the check-in record since there's no existing record for this attendee and faculty
	checkinErr := u.CheckinRepository.CreateCheckin(attendee.ID, staff.Faculty, staff.ID)
	if checkinErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create check-in",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok": true,
		"attendee": checkinModel.CheckinResponse{
			CheckedInAt: time.Now(),
			UserID:      attendee.UserID,
			Firstname:   attendee.Firstname,
			Surname:     attendee.Surname,
			TicketCode:  attendee.TicketCode,
			Faculty:     staff.Faculty,
		},
	})
}

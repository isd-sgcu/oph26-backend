package usecase

import (
	checkinModel "oph26-backend/internal/model/checkin"
	"oph26-backend/internal/repository"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type CheckinUsecase interface {
	CheckIn(c *fiber.Ctx) error
	GetCheckinStatus(c *fiber.Ctx) error
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
	reqID, _ := c.Locals("request_id").(string)
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
		slog.Error("checkin: failed to find attendee",
			"req_id", reqID, "ticket_code", req.TicketCode, "error", attendeeErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find attendee",
		})
	}

	if attendee == nil {
		slog.Warn("checkin: attendee not found",
			"req_id", reqID, "ticket_code", req.TicketCode, "staff_user_id", userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendee not found",
		})
	}

	staff, staffErr := u.StaffRepository.FindByUserID(userId)
	if staffErr != nil {
		slog.Error("checkin: failed to find staff",
			"req_id", reqID, "staff_user_id", userId, "error", staffErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find staff information",
		})
	}

	checkins, checkinsFindErr := u.CheckinRepository.FindCheckinByAttendeeAndFaculty(attendee.ID, staff.Faculty)
	if checkinsFindErr != nil {
		slog.Error("checkin: failed to query existing checkins",
			"req_id", reqID, "attendee_id", attendee.ID, "faculty", staff.Faculty, "error", checkinsFindErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check check-in status",
		})
	}

	if len(checkins) > 0 {
		firstCheckin := checkins[0]
		slog.Warn("checkin: already checked in",
			"req_id", reqID,
			"attendee_id", attendee.ID,
			"ticket_code", req.TicketCode,
			"faculty", staff.Faculty,
			"staff_user_id", userId,
			"existing_checkin_id", firstCheckin.ID,
			"existing_checked_in_at", firstCheckin.CheckedInAt,
			"duplicate_count", len(checkins),
		)
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

	checkinErr := u.CheckinRepository.CreateCheckin(attendee.ID, staff.Faculty, staff.ID)
	if checkinErr != nil {
		slog.Error("checkin: failed to create",
			"req_id", reqID, "attendee_id", attendee.ID, "faculty", staff.Faculty, "error", checkinErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create check-in",
		})
	}

	slog.Info("checkin: success",
		"req_id", reqID,
		"attendee_id", attendee.ID,
		"ticket_code", req.TicketCode,
		"faculty", staff.Faculty,
		"staff_user_id", userId,
	)

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

func (u *CheckinUsecaseImpl) GetCheckinStatus(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	attendee, err := u.AttendeeRepository.FindByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find attendee",
		})
	}
	if attendee == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendee not found",
		})
	}

	hasCheckedIn, statusErr := u.CheckinRepository.Checkinstatus(attendee.ID)
	if statusErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check check-in status",
		})
	}

	return c.JSON(fiber.Map{
		"status": hasCheckedIn,
	})
}

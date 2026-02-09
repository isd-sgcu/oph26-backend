package usecase

import (
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model"
	"oph26-backend/internal/model/attendee"
	"oph26-backend/internal/repository"
	"reflect"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AttendeesUsecase interface {
	GetMyAttendee(c *fiber.Ctx) error
	GetByAttendeeId(c *fiber.Ctx) error
	PutAttendeesUseCase(c *fiber.Ctx) error
}

type AttendeeUsecaseImpl struct {
	UserRepo     repository.UserRepository
	AttendeeRepo repository.AttendeeRepository
}

func NewAttendeeUsecase(userRepo repository.UserRepository, attendeeRepo repository.AttendeeRepository) AttendeesUsecase {
	return &AttendeeUsecaseImpl{
		UserRepo:     userRepo,
		AttendeeRepo: attendeeRepo,
	}
}

func (u *AttendeeUsecaseImpl) GetMyAttendee(c *fiber.Ctx) error {
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

	attendee, err := u.AttendeeRepo.FindByUserID(userID)
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

	return c.JSON(&model.AttendeeResponse{
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

func (u *AttendeeUsecaseImpl) GetByAttendeeId(c *fiber.Ctx) error {
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

	attendee, err := u.AttendeeRepo.FindByTicketCode(ticketCode)
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

	return c.JSON(&model.AttendeeStaffResponse{
		AttendeeResponse: model.AttendeeResponse{
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

func (u *AttendeeUsecaseImpl) PutAttendeesUseCase(c *fiber.Ctx) error {
	userEmail, ok := c.Locals("email").(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Could not assert email from JWT as string",
		})
	}
	userIdStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Could not assert user_id from JWT as string",
		})
	}
	userId, parseErr := uuid.Parse(userIdStr)
	if parseErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id",
		})
	}

	var reqBody attendee.PutAttendeesRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// User must exist and cannot be staff
	userFromDB, err := u.UserRepo.FindByEmail(userEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal DB error",
		})
	}
	if userFromDB == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attendee not found",
		})
	}
	if userFromDB.StaffId != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Update not allowed (not an attendee)",
		})
	}

	// No update to do if body is empty
	if (attendee.PutAttendeesRequest{}) == reqBody {
		return c.SendStatus(fiber.StatusNoContent)
	}

	// Validate enum fields
	if reqBody.InterestedFaculty != nil {
		arr := []string(*reqBody.InterestedFaculty)
		if !model.FacultiesAreValid(arr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body; unknown faculty",
			})
		}
	}
	if reqBody.Province != nil && !model.ProvinceIsValid(*reqBody.Province) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown province",
		})
	}
	if reqBody.NewsSourceSelected != nil {
		arr := []string(*reqBody.NewsSourceSelected)
		if !model.NewsSourcesAreValid(arr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body; unknown news source",
			})
		}
	}
	if reqBody.ObjectiveSelected != nil {
		arr := []string(*reqBody.ObjectiveSelected)
		if !model.ObjectivesAreValid(arr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body; unknown objective",
			})
		}
	}
	if reqBody.StudyLevel != nil && !model.StudyLevelIsValid(*reqBody.StudyLevel) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; unknown study_level",
		})
	}

	// If ObjectiveSelected = ["อื่น ๆ"], ObjectiveOther must have value
	tempOtherObj := pq.StringArray([]string{string(model.OtherObjective)})
	if reqBody.ObjectiveSelected == &tempOtherObj && reqBody.ObjectiveOther == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; objective_selected is 'อื่น ๆ', but objective_other is not provided",
		})
	}

	// If NewsSourceSelected = ["อื่น ๆ"], NewsSourcesOther must have value
	tempOtherNews := pq.StringArray([]string{string(model.OtherNewsSource)})
	if reqBody.NewsSourceSelected == &tempOtherNews && reqBody.NewsSourcesOther == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body; news_sources_selected is 'อื่น ๆ', but news_sources_other is not provided",
		})
	}

	// Map request body to attendee entity
	reqType := reflect.TypeOf(reqBody)
	updateStruct := entity.Attendee{}

	for i := 0; i < reqType.NumField(); i++ {
		switch reqType.Field(i).Name {
		case "Firstname":
			updateStruct.Firstname = *reqBody.Firstname

		case "Surname":
			updateStruct.Surname = *reqBody.Surname

		case "Age":
			updateStruct.Age = *reqBody.Age

		case "Province":
			updateStruct.Province = *reqBody.Province

		case "StudyLevel":
			updateStruct.StudyLevel = reqBody.StudyLevel

		case "SchoolName":
			updateStruct.SchoolName = reqBody.SchoolName

		case "NewsSourceSelected":
			updateStruct.NewsSourceSelected = *reqBody.NewsSourceSelected

		case "NewsSourcesOther":
			updateStruct.NewsSourcesOther = reqBody.NewsSourcesOther

		case "InterestedFaculty":
			updateStruct.InterestedFaculty = *reqBody.InterestedFaculty

		case "ObjectiveSelected":
			updateStruct.ObjectiveSelected = *reqBody.ObjectiveSelected

		case "ObjectiveOther":
			updateStruct.ObjectiveOther = reqBody.ObjectiveOther
		}
	}

	updateErr := u.AttendeeRepo.Update(&updateStruct, userId)
	if updateErr != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Internal DB error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok": true,
	})
}

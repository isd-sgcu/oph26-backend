package attendee

import (
	"oph26-backend/internal/model"
	"time"

	"github.com/google/uuid"
)

type AttendeeResponse struct {
	Age                           int        `json:"age"`
	AttendeeType                  string     `json:"attendee_type"`
	CertificateName               *string    `json:"certificate_name"`
	CheckedInAt                   *time.Time `json:"checked_in_at"`
	CheckinStaffID                *uuid.UUID `json:"checkin_staff_id"`
	CreatedAt                     time.Time  `json:"createdAt"`
	FavoriteWorkshops             []string   `json:"favorite_workshops"`
	Firstname                     string     `json:"firstname"`
	ID                            uuid.UUID  `json:"id"`
	InitialFirstInterestedFaculty string     `json:"initial_first_interested_faculty"`
	InterestedFaculty             []string   `json:"interested_faculty"`
	NewsSourcesOther              *string    `json:"news_sources_other"`
	NewsSourceSelected            []string   `json:"news_sources_selected"`
	ObjectiveOther                *string    `json:"objective_other"`
	ObjectiveSelected             []string   `json:"objective_selected"`
	Province                      string     `json:"province"`
	SchoolName                    *string    `json:"school_name"`
	StudyLevel                    *string    `json:"study_level"`
	Surname                       string     `json:"surname"`
	TicketCode                    string     `json:"ticket_code"`
	UpdatedAt                     time.Time  `json:"updatedAt"`
	UserID                        uuid.UUID  `json:"user_id"`
}

type AttendeeStaffResponse struct {
	AttendeeResponse
	CheckinStaff *model.StaffResponse `json:"checkin_staff"`
}

type GetFavoriteWorkshopResponse struct {
	FavoriteWorkshop []string `json:"favorite_workshop"`
}

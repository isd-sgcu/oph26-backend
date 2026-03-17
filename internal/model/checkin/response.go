package checkin

import (
	"oph26-backend/internal/model"
	"oph26-backend/internal/model/attendee"
	"time"

	"github.com/google/uuid"
)

type CheckinEntry struct {
	ID          uuid.UUID `json:"id"`
	AttendeeID  uuid.UUID `json:"attendee_id"`
	Faculty     string    `json:"faculty"`
	StaffID     uuid.UUID `json:"staff_id"`
	CheckedInAt string    `json:"checked_in_at"`
}

type CheckinEntryConflictResponse struct {
	Error string `json:"error"`
	CheckinEntry
	Attendee attendee.AttendeeResponse `json:"attendee"`
	Staff    model.StaffResponse       `json:"staff"`
}
type CheckinResponse struct {
	CheckedInAt time.Time `json:"checked_in_at"`
	UserID      uuid.UUID `json:"user_id"`
	Firstname   string    `json:"firstname"`
	Surname     string    `json:"surname"`
	TicketCode  string    `json:"ticket_code"`
}

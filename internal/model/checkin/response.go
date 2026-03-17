package checkin

import (
	"oph26-backend/internal/model"
	"oph26-backend/internal/model/attendee"

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

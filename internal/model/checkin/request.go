package checkin

type CheckiAttendeeRequest struct {
	TicketCode string `json:"ticket_code" validate:"required"`
}

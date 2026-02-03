package user

import (
	"time"

	"github.com/google/uuid"
)

type GetAllUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	AttendeeId *uuid.UUID `json:"attendee_id,omitempty"`
	StaffId    *uuid.UUID `json:"staff_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	Email      string
	Role       string
	AttendeeId *uuid.UUID
	StaffId    *uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUser(email, role string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type StaffResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"userId"`
	Cuid      string     `json:"cuid"`
	Firstname string     `json:"firstname"`
	Surname   string     `json:"surname"`
	Nickname  string     `json:"nickname"`
	Phone     string     `json:"phone"`
	Year      string     `json:"year"`
	Email     string     `json:"email"`
	Faculty   string     `json:"faculty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

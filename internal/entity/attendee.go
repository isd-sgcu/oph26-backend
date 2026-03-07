package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/lib/pq"
)

type Attendee struct {
	ID                            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID                        uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex"`
	User                          User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Firstname                     string         `gorm:"not null"`
	Surname                       string         `gorm:"not null"`
	AttendeeType                  string         `gorm:"type:text;not null;check:attendee_type IN ('elementaryschool','highschool','parent','educationstaff','other')"`
	Age                           int            `gorm:"not null"`
	Province                      string         `gorm:"not null"`
	StudyLevel                    *string        `gorm:"type:text"`
	SchoolName                    *string        `gorm:"type:text"`
	NewsSourceSelected            pq.StringArray `gorm:"type:text[]"`
	NewsSourcesOther              *string        `gorm:"type:text"`
	InitialFirstInterestedFaculty string         `gorm:"not null"`
	InterestedFaculty             pq.StringArray `gorm:"type:text[];not null"`
	ObjectiveSelected             pq.StringArray `gorm:"type:text[]"`
	ObjectiveOther                *string        `gorm:"type:text"`
	TicketCode                    string         `gorm:"type:char(7);not null;uniqueIndex"`
	MyPiece                       *MyPiece
	CertificateName               *string    `gorm:"type:text"`
	CheckinStaffID                *uuid.UUID `gorm:"type:uuid"`
	CheckinStaff                  *Staff     `gorm:"foreignKey:CheckinStaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CheckedInAt                   *time.Time
	FavoriteWorkshops             pq.StringArray `gorm:"type:text[]"`
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
}

package entity

import (
	"database/sql/driver"
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
	DateOfBirth                   time.Time      `gorm:"type:date;not null"`
	Province                      string         `gorm:"not null"`
	District                      string         `gorm:"not null"`
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
	FavoriteWorkshops             StringSet `gorm:"type:text[]"`
	TransportationMethod          string
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
}

type StringSet map[string]struct{}

func (s *StringSet) Scan(value any) error {
	if value == nil {
		*s = make(StringSet)
		return nil
	}
	var arr pq.StringArray
	if err := arr.Scan(value); err != nil {
		return err
	}
	set := make(StringSet)
	for _, item := range arr {
		set[item] = struct{}{}
	}
	*s = set
	return nil
}

func (s StringSet) Value() (driver.Value, error) {
	if s == nil {
		return pq.StringArray{}.Value()
	}
	arr := make([]string, 0, len(s))
	for item := range s {
		arr = append(arr, item)
	}
	return pq.StringArray(arr).Value()
}

func (s StringSet) ToSlice() []string {
	result := make([]string, 0, len(s))
	for item := range s {
		result = append(result, item)
	}
	return result
}

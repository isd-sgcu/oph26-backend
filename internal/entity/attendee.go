package entity

import (
	"database/sql/driver"
	"time"

	"oph26-backend/internal/model"

	"github.com/google/uuid"

	"github.com/lib/pq"
)

type Attendee struct {
	ID                            uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID                        uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex"`
	User                          User              `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Firstname                     string            `gorm:"not null"`
	Surname                       string            `gorm:"not null"`
	AttendeeType                  string            `gorm:"type:text;not null;check:attendee_type IN ('student','parent','educationstaff','other')"`
	DateOfBirth                   time.Time         `gorm:"type:date;not null"`
	Province                      string            `gorm:"not null"`
	District                      string            `gorm:"not null"`
	StudyLevel                    *model.StudyLevel `gorm:"type:text"`
	SchoolName                    *string           `gorm:"type:text"`
	NewsSourceSelected            pq.StringArray    `gorm:"type:text[]"`
	NewsSourcesOther              *string           `gorm:"type:text"`
	InitialFirstInterestedFaculty model.Faculty     `gorm:"type:text"`
	InterestedFaculty             FacultyList       `gorm:"type:text[]"`
	ObjectiveSelected             pq.StringArray    `gorm:"type:text[]"`
	ObjectiveOther                *string           `gorm:"type:text"`
	TicketCode                    string            `gorm:"type:char(7);not null;uniqueIndex"`
	MyPiece                       *MyPiece
	CertificateName               *string    `gorm:"type:text"`
	CheckinStaffID                *uuid.UUID `gorm:"type:uuid"`
	CheckinStaff                  *Staff     `gorm:"foreignKey:CheckinStaffID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CheckedInAt                   *time.Time
	FavoriteWorkshops             StringSet `gorm:"type:text[]"`
	TransportationMethod          string
	Rank                          int `gorm:"default:-1"`
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

type FacultyList []model.Faculty

func (f FacultyList) Value() (driver.Value, error) {
	strs := make([]string, len(f))
	for i, fac := range f {
		strs[i] = string(fac)
	}
	return pq.Array(strs).Value()
}

func (f *FacultyList) Scan(value interface{}) error {
	var strs pq.StringArray
	if err := strs.Scan(value); err != nil {
		return err
	}
	*f = make(FacultyList, len(strs))
	for i, s := range strs {
		(*f)[i] = model.Faculty(s)
	}
	return nil
}

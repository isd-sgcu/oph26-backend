package repository

import (
	"oph26-backend/internal/entity"

	"gorm.io/gorm"
)

type StatsRepository interface {
	CountAttendees() (int64, error)
	CountAttendeesGroupedByType() (map[string]int64, error)
	CountUniqueAttendeesCheckinsGroupedByDateAndType() (map[string]map[string]int64, error)
	CountCheckins() (int64, error)
	CountCheckinsGroupedByDate() (map[string]int64, error)
	CountUniqueAttendeesCheckinsGroupedByDate() (map[string]int64, error)
	CountAvailablePiecesGroupedByFaculty() (map[string]int64, error)
	// Umm should we implement these?
	// CountAttendeeWithFullyCollectedPieces() (int64, error)
}

type StatsRepositoryImpl struct {
	DB *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &StatsRepositoryImpl{DB: db}
}

func (r *StatsRepositoryImpl) CountAttendees() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Attendee{}).Count(&count).Error
	return count, err
}

func (r *StatsRepositoryImpl) CountAttendeesGroupedByType() (map[string]int64, error) {
	type Result struct {
		Type  string
		Count int64
	}
	var results []Result
	err := r.DB.Model(&entity.Attendee{}).
		Select("attendee_type AS type, COUNT(*) AS count").
		Group("attendee_type").
		Order("attendee_type").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, res := range results {
		counts[res.Type] = res.Count
	}
	return counts, nil
}

func (r *StatsRepositoryImpl) CountUniqueAttendeesCheckinsGroupedByDateAndType() (map[string]map[string]int64, error) {
	type Result struct {
		Date         string
		AttendeeType string
		Count        int64
	}

	var results []Result
	err := r.DB.Model(&entity.Checkin{}).
		Select("DATE(checkins.checked_in_at) AS date, attendees.attendee_type AS attendee_type, COUNT(DISTINCT checkins.attendee_id) AS count").
		Joins("JOIN attendees ON attendees.id = checkins.attendee_id").
		Group("date, attendees.attendee_type").
		Order("date, attendees.attendee_type").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]map[string]int64)
	for _, res := range results {
		if _, ok := counts[res.Date]; !ok {
			counts[res.Date] = make(map[string]int64)
		}
		counts[res.Date][res.AttendeeType] = res.Count
	}

	return counts, nil
}

func (r *StatsRepositoryImpl) CountCheckins() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Checkin{}).Count(&count).Error
	return count, err
}

func (r *StatsRepositoryImpl) CountCheckinsGroupedByDate() (map[string]int64, error) {
	type Result struct {
		Date  string
		Count int64
	}

	var results []Result
	err := r.DB.Model(&entity.Checkin{}).
		Select("DATE(checked_in_at) AS date, COUNT(*) AS count").
		Group("date").
		Order("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, res := range results {
		counts[res.Date] = res.Count
	}
	return counts, nil
}

func (r *StatsRepositoryImpl) CountUniqueAttendeesCheckinsGroupedByDate() (map[string]int64, error) {
	type Result struct {
		Date  string
		Count int64
	}
	var results []Result
	err := r.DB.Model(&entity.Checkin{}).
		Select("DATE(checked_in_at) AS date, COUNT(DISTINCT attendee_id) AS count").
		Group("date").
		Order("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, res := range results {
		counts[res.Date] = res.Count
	}
	return counts, nil
}

func (r *StatsRepositoryImpl) CountAvailablePiecesGroupedByFaculty() (map[string]int64, error) {
	type Result struct {
		Faculty string
		Count   int64
	}

	var results []Result
	err := r.DB.Model(&entity.MyPiece{}).
		Select("attendees.initial_first_interested_faculty AS faculty, COUNT(*) AS count").
		Joins("JOIN attendees ON attendees.id = my_pieces.attendee_id").
		Group("faculty").
		Order("faculty").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, res := range results {
		counts[res.Faculty] = res.Count
	}
	return counts, nil
}

package metrics

import (
	"oph26-backend/internal/repository"

	"github.com/prometheus/client_golang/prometheus"
)

type AttendeeMetrics struct {
	statsRepo                            repository.StatsRepository
	attendeesByType                      *prometheus.GaugeVec
	knownTypeLabels                      []string
	attendeesTotal                       prometheus.Gauge
	uniqueAttendeesCheckinsByDate        *prometheus.GaugeVec
	availablePiecesByFaculty             *prometheus.GaugeVec
	uniqueAttendeesCheckinsByDateAndType *prometheus.GaugeVec
	checkinsTotal                        prometheus.Gauge
	checkinsByDate                       *prometheus.GaugeVec
}

func NewAttendeeMetrics(statsRepo repository.StatsRepository) *AttendeeMetrics {
	attendeesByType := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_attendees_by_type",
			Help: "Current number of attendees grouped by attendee type",
		},
		[]string{"attendee_type"},
	)

	attendeesTotal := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cuoph26_attendees_total",
			Help: "Current number of attendees",
		},
	)

	uniqueAttendeesCheckinsByDate := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_unique_attendees_checkins_by_date",
			Help: "Unique attendees checked in, grouped by checkin date",
		},
		[]string{"checkin_date"},
	)

	availablePiecesByFaculty := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_available_pieces_by_faculty",
			Help: "Current available pieces grouped by faculty",
		},
		[]string{"faculty"},
	)

	uniqueAttendeesCheckinsByDateAndType := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_unique_attendees_checkins_by_date_and_type",
			Help: "Unique attendees checked in, grouped by checkin date and attendee type",
		},
		[]string{"checkin_date", "attendee_type"},
	)

	checkinsTotal := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cuoph26_checkins_total",
			Help: "Current total number of checkins",
		},
	)

	checkinsByDate := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_checkins_by_date",
			Help: "Current total checkins grouped by checkin date",
		},
		[]string{"checkin_date"},
	)

	prometheus.MustRegister(
		attendeesByType,
		attendeesTotal,
		uniqueAttendeesCheckinsByDate,
		availablePiecesByFaculty,
		uniqueAttendeesCheckinsByDateAndType,
		checkinsTotal,
		checkinsByDate,
	)

	return &AttendeeMetrics{
		statsRepo:                            statsRepo,
		attendeesByType:                      attendeesByType,
		knownTypeLabels:                      []string{"student", "parent", "educationstaff", "other"},
		attendeesTotal:                       attendeesTotal,
		uniqueAttendeesCheckinsByDate:        uniqueAttendeesCheckinsByDate,
		availablePiecesByFaculty:             availablePiecesByFaculty,
		uniqueAttendeesCheckinsByDateAndType: uniqueAttendeesCheckinsByDateAndType,
		checkinsTotal:                        checkinsTotal,
		checkinsByDate:                       checkinsByDate,
	}
}

func (m *AttendeeMetrics) Refresh() error {
	attendeesCount, err := m.statsRepo.CountAttendees()
	if err != nil {
		return err
	}
	m.attendeesTotal.Set(float64(attendeesCount))

	countsByType, err := m.statsRepo.CountAttendeesGroupedByType()
	if err != nil {
		return err
	}

	for _, attendeeType := range m.knownTypeLabels {
		m.attendeesByType.WithLabelValues(attendeeType).Set(0)
	}

	for attendeeType, count := range countsByType {
		m.attendeesByType.WithLabelValues(attendeeType).Set(float64(count))
	}

	uniqueByDate, err := m.statsRepo.CountUniqueAttendeesCheckinsGroupedByDate()
	if err != nil {
		return err
	}
	m.uniqueAttendeesCheckinsByDate.Reset()
	for checkinDate, count := range uniqueByDate {
		m.uniqueAttendeesCheckinsByDate.WithLabelValues(checkinDate).Set(float64(count))
	}

	piecesByFaculty, err := m.statsRepo.CountAvailablePiecesGroupedByFaculty()
	if err != nil {
		return err
	}
	m.availablePiecesByFaculty.Reset()
	for faculty, count := range piecesByFaculty {
		m.availablePiecesByFaculty.WithLabelValues(faculty).Set(float64(count))
	}

	uniqueByDateAndType, err := m.statsRepo.CountUniqueAttendeesCheckinsGroupedByDateAndType()
	if err != nil {
		return err
	}
	m.uniqueAttendeesCheckinsByDateAndType.Reset()
	for checkinDate, countsByAttendeeType := range uniqueByDateAndType {
		for attendeeType, count := range countsByAttendeeType {
			m.uniqueAttendeesCheckinsByDateAndType.WithLabelValues(checkinDate, attendeeType).Set(float64(count))
		}
	}

	checkinsCount, err := m.statsRepo.CountCheckins()
	if err != nil {
		return err
	}
	m.checkinsTotal.Set(float64(checkinsCount))

	checkinsByDate, err := m.statsRepo.CountCheckinsGroupedByDate()
	if err != nil {
		return err
	}
	m.checkinsByDate.Reset()
	for checkinDate, count := range checkinsByDate {
		m.checkinsByDate.WithLabelValues(checkinDate).Set(float64(count))
	}

	return nil
}

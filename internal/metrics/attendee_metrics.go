package metrics

import (
	"oph26-backend/internal/repository"
	"time"

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
	checkinsByFaculty                    *prometheus.GaugeVec
	checkinsByStaff                      *prometheus.GaugeVec
	checkinsByHourAndFaculty             *prometheus.GaugeVec
	uniqueAttendeesCheckedInToday        prometheus.Gauge
	duplicateCheckinsToday               prometheus.Gauge
	attendeesWithCompletedPieces         prometheus.Gauge
	totalCollectedPieces                 prometheus.Gauge
	myPiecesGroupedByFaculty             *prometheus.GaugeVec
	maxPiecesCollectedByOneAttendee      prometheus.Gauge
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

	checkinsByFaculty := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_checkins_by_faculty",
			Help: "Current total checkins grouped by faculty",
		},
		[]string{"faculty"},
	)

	checkinsByStaff := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_checkins_by_staff",
			Help: "Current total checkins grouped by staff",
		},
		[]string{"staff"},
	)

	checkinsByHourAndFaculty := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_checkins_by_hour_and_faculty",
			Help: "Current total checkins grouped by hour bucket and faculty",
		},
		[]string{"hour_bucket", "faculty"},
	)

	uniqueAttendeesCheckedInToday := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cuoph26_unique_attendees_checked_in_today",
			Help: "Unique attendees checked in today",
		},
	)

	duplicateCheckinsToday := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cuoph26_duplicate_checkins_today",
			Help: "Duplicate checkins today (total checkins today - unique attendees checked in today)",
		},
	)

	attendeesWithCompletedPieces := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cuoph26_attendees_with_completed_pieces",
			Help: "Number of attendees who have completed all pieces",
		},
	)

	totalCollectedPieces := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cuoph26_total_collected_pieces",
			Help: "Total number of pieces collected by all attendees",
		},
	)

	myPiecesGroupedByFaculty := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cuoph26_my_pieces_by_faculty",
			Help: "Total pieces (MyPieces) available grouped by faculty",
		},
		[]string{"faculty"},
	)

	maxPiecesCollectedByOneAttendee := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cuoph26_max_pieces_collected_by_one_attendee",
			Help: "Maximum number of pieces collected by a single attendee",
		},
	)

	prometheus.MustRegister(
		attendeesByType,
		attendeesTotal,
		uniqueAttendeesCheckinsByDate,
		availablePiecesByFaculty,
		uniqueAttendeesCheckinsByDateAndType,
		checkinsTotal,
		checkinsByDate,
		checkinsByFaculty,
		checkinsByStaff,
		checkinsByHourAndFaculty,
		uniqueAttendeesCheckedInToday,
		duplicateCheckinsToday,
		attendeesWithCompletedPieces,
		totalCollectedPieces,
		myPiecesGroupedByFaculty,
		maxPiecesCollectedByOneAttendee,
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
		checkinsByFaculty:                    checkinsByFaculty,
		checkinsByStaff:                      checkinsByStaff,
		checkinsByHourAndFaculty:             checkinsByHourAndFaculty,
		uniqueAttendeesCheckedInToday:        uniqueAttendeesCheckedInToday,
		duplicateCheckinsToday:               duplicateCheckinsToday,
		attendeesWithCompletedPieces:         attendeesWithCompletedPieces,
		totalCollectedPieces:                 totalCollectedPieces,
		myPiecesGroupedByFaculty:             myPiecesGroupedByFaculty,
		maxPiecesCollectedByOneAttendee:      maxPiecesCollectedByOneAttendee,
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

	checkinsByFaculty, err := m.statsRepo.CountCheckinsGroupedByFaculty()
	if err != nil {
		return err
	}
	m.checkinsByFaculty.Reset()
	for faculty, count := range checkinsByFaculty {
		m.checkinsByFaculty.WithLabelValues(faculty).Set(float64(count))
	}

	checkinsByStaff, err := m.statsRepo.CountCheckinsGroupedByStaff()
	if err != nil {
		return err
	}
	m.checkinsByStaff.Reset()
	for staff, count := range checkinsByStaff {
		m.checkinsByStaff.WithLabelValues(staff).Set(float64(count))
	}

	checkinsByHourAndFaculty, err := m.statsRepo.CountCheckinsGroupedByHourAndFacultySince(time.Now().Add(-24 * time.Hour))
	if err != nil {
		return err
	}
	m.checkinsByHourAndFaculty.Reset()
	for hourBucket, facultyCounts := range checkinsByHourAndFaculty {
		for faculty, count := range facultyCounts {
			m.checkinsByHourAndFaculty.WithLabelValues(hourBucket, faculty).Set(float64(count))
		}
	}

	uniqueToday, err := m.statsRepo.CountUniqueAttendeesCheckedInToday()
	if err != nil {
		return err
	}
	m.uniqueAttendeesCheckedInToday.Set(float64(uniqueToday))

	duplicateToday, err := m.statsRepo.CountDuplicateCheckinsToday()
	if err != nil {
		return err
	}
	m.duplicateCheckinsToday.Set(float64(duplicateToday))

	attendeesWithCompletedPieces, err := m.statsRepo.CountAttendeeWithCompletedPieces()
	if err != nil {
		return err
	}
	m.attendeesWithCompletedPieces.Set(float64(attendeesWithCompletedPieces))

	totalCollected, err := m.statsRepo.CountTotalCollectedPieces()
	if err != nil {
		return err
	}
	m.totalCollectedPieces.Set(float64(totalCollected))

	myPiecesGroupedByFacultyData, err := m.statsRepo.CountMyPiecesGroupedByFaculty()
	if err != nil {
		return err
	}
	m.myPiecesGroupedByFaculty.Reset()
	for faculty, count := range myPiecesGroupedByFacultyData {
		m.myPiecesGroupedByFaculty.WithLabelValues(faculty).Set(float64(count))
	}

	maxPiecesCollected, err := m.statsRepo.GetMaxPiecesCollectedByOneAttendee()
	if err != nil {
		return err
	}
	m.maxPiecesCollectedByOneAttendee.Set(float64(maxPiecesCollected))

	return nil
}

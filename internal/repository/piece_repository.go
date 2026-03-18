package repository

import (
	"errors"
	"oph26-backend/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PieceRepository interface {
	FindAttendeeByUserID(userID uuid.UUID) (*entity.Attendee, error)
	FindMyPieceByAttendeeID(attendeeID uuid.UUID) (*entity.MyPiece, error)
	RefreshMyPiece(piece *entity.MyPiece, newCode string) error
	FindMyPieceByCode(pieceCode string) (*entity.MyPiece, error)
	FindCollectedPiecesByAttendeeID(attendeeID uuid.UUID) ([]entity.CollectedPiece, error)
	FindCollectedPiece(attendeeID uuid.UUID, pieceID uuid.UUID) (*entity.CollectedPiece, error)
	CreateCollectedPiece(cp *entity.CollectedPiece) error
	CountCollectedByFaculty(attendeeID uuid.UUID) (map[string]int, error)
	CountTop1ThresholdByFaculty() (map[string]int, error)
}

type PieceRepositoryImpl struct {
	DB *gorm.DB
}

func NewPieceRepository(db *gorm.DB) PieceRepository {
	return &PieceRepositoryImpl{DB: db}
}

func (r *PieceRepositoryImpl) FindAttendeeByUserID(userID uuid.UUID) (*entity.Attendee, error) {
	var attendee entity.Attendee
	if err := r.DB.Where("user_id = ?", userID).First(&attendee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attendee, nil
}

func (r *PieceRepositoryImpl) FindMyPieceByAttendeeID(attendeeID uuid.UUID) (*entity.MyPiece, error) {
	var piece entity.MyPiece
	if err := r.DB.Where("attendee_id = ?", attendeeID).Order("created_at DESC").First(&piece).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &piece, nil
}

func (r *PieceRepositoryImpl) RefreshMyPiece(piece *entity.MyPiece, newCode string) error {
	newExpireDate := time.Now().Add(24 * time.Hour)
	if err := r.DB.Model(piece).Updates(map[string]interface{}{
		"piece_code":  newCode,
		"expire_date": newExpireDate,
	}).Error; err != nil {
		return err
	}

	piece.PieceCode = newCode
	piece.ExpireDate = newExpireDate
	return nil
}
  
func (r *PieceRepositoryImpl) FindMyPieceByCode(pieceCode string) (*entity.MyPiece, error) {
	var piece entity.MyPiece
	if err := r.DB.Preload("Attendee").Where("piece_code = ?", pieceCode).First(&piece).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &piece, nil
}

func (r *PieceRepositoryImpl) FindCollectedPiece(attendeeID uuid.UUID, pieceID uuid.UUID) (*entity.CollectedPiece, error) {
	var cp entity.CollectedPiece
	if err := r.DB.Where("attendee_id = ? AND piece_id = ?", attendeeID, pieceID).First(&cp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cp, nil
}

func (r *PieceRepositoryImpl) CreateCollectedPiece(cp *entity.CollectedPiece) error {
	return r.DB.Create(cp).Error
}

func (r *PieceRepositoryImpl) FindCollectedPiecesByAttendeeID(attendeeID uuid.UUID) ([]entity.CollectedPiece, error) {
	var pieces []entity.CollectedPiece
	if err := r.DB.
		Preload("MyPiece").
		Preload("MyPiece.Attendee").
		Where("attendee_id = ?", attendeeID).
		Find(&pieces).Error; err != nil {
		return nil, err
	}
	return pieces, nil
}

func (r *PieceRepositoryImpl) CountCollectedByFaculty(attendeeID uuid.UUID) (map[string]int, error) {
	type result struct {
		Faculty string
		Count   int
	}
	var results []result

	err := r.DB.
		Table("collected_pieces").
		Select("attendees.initial_first_interested_faculty AS faculty, COUNT(*) AS count").
		Joins("JOIN my_pieces ON my_pieces.id = collected_pieces.piece_id").
		Joins("JOIN attendees ON attendees.id = my_pieces.attendee_id").
		Where("collected_pieces.attendee_id = ?", attendeeID).
		Group("attendees.initial_first_interested_faculty").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)
	for _, r := range results {
		m[r.Faculty] = r.Count
	}
	return m, nil
}

func (r *PieceRepositoryImpl) CountTop1ThresholdByFaculty() (map[string]int, error) {
	type result struct {
		Faculty   string
		Threshold int
	}
	var results []result

	err := r.DB.Raw(`
		SELECT faculty, COALESCE(percentile_disc(0.99) WITHIN GROUP (ORDER BY cnt), 0)::int AS threshold
		FROM (
			SELECT attendees.initial_first_interested_faculty AS faculty,
			       collected_pieces.attendee_id,
			       COUNT(*) AS cnt
			FROM collected_pieces
			JOIN my_pieces ON my_pieces.id = collected_pieces.piece_id
			JOIN attendees ON attendees.id = my_pieces.attendee_id
			GROUP BY attendees.initial_first_interested_faculty, collected_pieces.attendee_id
		) sub
		GROUP BY faculty
	`).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)
	for _, r := range results {
		m[r.Faculty] = r.Threshold
	}
	return m, nil
}

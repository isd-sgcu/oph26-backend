package piece

import (
	"oph26-backend/internal/model"
	"time"

	"github.com/google/uuid"
)

type MyPieceResponse struct {
	ID         uuid.UUID     `json:"id"`
	UserID     uuid.UUID     `json:"user_id"`
	PieceCode  string        `json:"piece_code"`
	ExpireDate time.Time     `json:"expire_date"`
	Faculty    model.Faculty `json:"faculty"`
}

type FriendPieceResponse struct {
	ID          uuid.UUID     `json:"id"`
	UserID      uuid.UUID     `json:"user_id"`
	Faculty     model.Faculty `json:"faculty"`
	CollectedAt *time.Time    `json:"collected_at,omitempty"`
}

type FacultyStats struct {
	Count  int  `json:"count"`
	IsTop1 bool `json:"is_top_1"`
}

type CollectedPiecesStats struct {
	TotalCollected     int                     `json:"total_collected"`
	CollectedByFaculty map[string]FacultyStats `json:"collected_by_faculty"`
	SameMissingCount   map[int]float64         `json:"same_missing_count"`
	Rank               int                     `json:"rank"`
}

type CollectedPiecesResponse struct {
	CollectedPieces []FriendPieceResponse `json:"collected_pieces"`
	Stats           CollectedPiecesStats  `json:"stats"`
}

type CollectPieceRequest struct {
	PieceCode string `json:"piece_code" validate:"required"`
}

type CollectPieceResponse struct {
	Ok             bool                `json:"ok"`
	CollectedPiece FriendPieceResponse `json:"collected_piece"`
}

package piece

import (
	"time"

	"github.com/google/uuid"
)

type MyPieceResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	PieceCode  string    `json:"piece_code"`
	ExpireDate time.Time `json:"expire_date"`
	Faculty    string    `json:"faculty"`
}

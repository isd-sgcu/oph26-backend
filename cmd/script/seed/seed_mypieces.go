package main

import (
	"log"
	"oph26-backend/internal/entity"
	"time"

	"gorm.io/gorm"
)

// seedMyPieces seeds one piece per highschool attendee (attendees[0] and [1]).
// attendees[0] gets an already-expired piece to demonstrate the auto-refresh flow.
func seedMyPieces(db *gorm.DB, attendees []entity.Attendee) []entity.MyPiece {
	pieces := []entity.MyPiece{
		{
			AttendeeID: attendees[0].ID,
			PieceCode:  "EXP001",
			ExpireDate: time.Now().Add(-1 * time.Hour), // expired — triggers refresh on next GET
		},
		{
			AttendeeID: attendees[1].ID,
			PieceCode:  "ACTIVE1",
			ExpireDate: time.Now().Add(24 * time.Hour),
		},
	}

	for i := range pieces {
		if err := db.Where("piece_code = ?", pieces[i].PieceCode).FirstOrCreate(&pieces[i]).Error; err != nil {
			log.Fatalf("Failed to seed piece %s: %v", pieces[i].PieceCode, err)
		}
	}

	log.Printf("Seeded %d my-pieces", len(pieces))
	return pieces
}

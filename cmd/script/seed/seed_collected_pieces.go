package main

import (
	"log"
	"oph26-backend/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// seedCollectedPieces seeds cross-collection between attendees[0] and attendees[1].
func seedCollectedPieces(db *gorm.DB, attendees []entity.Attendee, pieces []entity.MyPiece) {
	collectedPieces := []entity.CollectedPiece{
		{
			AttendeeID:  attendees[0].ID,
			PieceID:     pieces[1].ID, // attendee[0] collected attendee[1]'s piece
			CollectedAt: time.Now(),
		},
		{
			AttendeeID:  attendees[1].ID,
			PieceID:     pieces[0].ID, // attendee[1] collected attendee[0]'s piece
			CollectedAt: time.Now(),
		},
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&collectedPieces).Error; err != nil {
		log.Fatalf("Failed to seed collected pieces: %v", err)
	}

	log.Printf("Seeded %d collected pieces", len(collectedPieces))
}

package main

import (
	"log"
	"oph26-backend/internal/entity"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// seedScores seeds a zeroed 20-element score array for each of the 3 attendee users.
func seedScores(db *gorm.DB, users []entity.User) {
	scores := []entity.Score{
		{UserID: users[0].ID, Count: pq.Int32Array(make([]int32, 20))},
		{UserID: users[1].ID, Count: pq.Int32Array(make([]int32, 20))},
		{UserID: users[2].ID, Count: pq.Int32Array(make([]int32, 20))},
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&scores).Error; err != nil {
		log.Fatalf("Failed to seed scores: %v", err)
	}

	log.Printf("Seeded %d scores", len(scores))
}

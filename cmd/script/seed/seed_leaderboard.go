package main

import (
	"log"
	"oph26-backend/internal/entity"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// seedLeaderboard seeds a zeroed 20-element bool array for each of the 3 attendee users.
func seedLeaderboard(db *gorm.DB, users []entity.User) {
	entries := []entity.Leaderboard{
		{UserID: users[0].ID, IsTop: pq.BoolArray(make([]bool, 20))},
		{UserID: users[1].ID, IsTop: pq.BoolArray(make([]bool, 20))},
		{UserID: users[2].ID, IsTop: pq.BoolArray(make([]bool, 20))},
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entries).Error; err != nil {
		log.Fatalf("Failed to seed leaderboard: %v", err)
	}

	log.Printf("Seeded %d leaderboard entries", len(entries))
}

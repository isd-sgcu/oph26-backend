package main

import (
	"log"
	"oph26-backend/internal/config"
	"oph26-backend/internal/initializer"
)

func init() {
	initializer.LoadEnvVariables()
}

func main() {
	cfg := config.LoadEnv()
	config.InitDB(cfg)

	db := config.DB
	log.Println("Starting seed...")

	// Order matters: respect FK dependencies
	users := seedUsers(db)
	staffs := seedStaff(db, users)
	attendees := seedAttendees(db, users, staffs)
	pieces := seedMyPieces(db, attendees)
	seedCollectedPieces(db, attendees, pieces)
	seedScores(db, users)
	seedLeaderboard(db, users)

	// Note: RefreshToken is runtime-generated (skipped)
	// Note: Questionnaire.UserID is type string referencing uuid — seed manually if needed

	log.Println("Seed completed successfully!")
}

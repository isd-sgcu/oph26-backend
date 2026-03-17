package main

import (
	"fmt"
	"log"
	"oph26-backend/internal/config"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/initializer"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	initializer.LoadEnvVariables()
}

func main() {
	cfg := config.LoadEnv()

	db, err := gorm.Open(postgres.Open(cfg.DataBaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Dropping all tables...")
	err = db.Migrator().DropTable(
		&entity.CollectedPiece{},
		&entity.MyPiece{},
		&entity.Score{},
		&entity.Leaderboard{},
		&entity.Questionnaire{},
		&entity.RefreshToken{},
		&entity.Attendee{},
		&entity.Staff{},
		&entity.User{},
	)
	if err != nil {
		log.Fatal("Failed to drop tables:", err)
	}

	fmt.Println("Re-creating extension and migrating...")
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Score{},
		&entity.Leaderboard{},
		&entity.Staff{},
		&entity.Attendee{},
		&entity.MyPiece{},
		&entity.CollectedPiece{},
		&entity.RefreshToken{},
		&entity.Questionnaire{},
	)
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	fmt.Println("DB reset done.")

	fmt.Println("Seeding data...")

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

	log.Println("Seeding completed successfully!")
}

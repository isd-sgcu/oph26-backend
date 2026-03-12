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

func main() {
	initializer.LoadEnvVariables()
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
}

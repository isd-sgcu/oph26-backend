package config

import (
	"log"
	"oph26-backend/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *Config) {
	var err error
	dsn := cfg.DataBaseURL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	// Enable uuid-ossp extension for UUID generation
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatal("Failed to create uuid-ossp extension:", err)
		return
	}
	err = DB.AutoMigrate(&entity.User{}, &entity.Score{}, &entity.Leaderboard{}, &entity.Staff{}, &entity.Attendee{}, &entity.MyPiece{}, &entity.CollectedPiece{}, &entity.RefreshToken{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return
	}
	log.Println("Database connection established and migrated")
}

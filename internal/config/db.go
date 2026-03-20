package config

import (
	"log"
	"oph26-backend/internal/entity"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *Config) {
	var err error
	dsn := cfg.DataBaseURL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
		return
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(3 * time.Minute)
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatal("Failed to create uuid-ossp extension:", err)
		return
	}
	err = DB.AutoMigrate(&entity.User{}, &entity.Score{}, &entity.Leaderboard{}, &entity.Staff{}, &entity.Attendee{}, &entity.MyPiece{}, &entity.CollectedPiece{}, &entity.RefreshToken{}, &entity.Questionnaire{}, &entity.Checkin{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return
	}
	log.Println("Database connection established and migrated")
}

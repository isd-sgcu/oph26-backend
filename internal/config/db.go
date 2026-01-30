package config

import (
	"fmt"
	"oph26-backend/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *Config) {
	var err error
	dsn := cfg.DataBaseURL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	// Enable uuid-ossp extension for UUID generation
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		fmt.Println("Failed to create uuid-ossp extension:", err)
		return
	}
	err = DB.AutoMigrate(&model.User{}, &model.Score{}, &model.Leaderboard{}, &model.Staff{}, &model.MyPiece{})
	if err != nil {
		fmt.Println("Failed to migrate database:", err)
		return
	}
	fmt.Println("Database connection established and migrated")
}

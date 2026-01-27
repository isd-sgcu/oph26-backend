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
	err = DB.AutoMigrate(&model.Staff{})
	if err != nil {
		fmt.Println("Failed to migrate database:", err)
		return
	}
	fmt.Println("Database connection established and migrated")
}

package main

import (
	"log"
	"oph26-backend/internal/entity"

	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) []entity.User {
	users := []entity.User{
		{Email: "attendee1@test.com", Role: "user"},
		{Email: "attendee2@test.com", Role: "user"},
		{Email: "attendee3@test.com", Role: "user"},
		{Email: "staff1@chula.ac.th", Role: "staff"},
		{Email: "staff2@chula.ac.th", Role: "staff"},
	}

	for i := range users {
		if err := db.Where("email = ?", users[i].Email).FirstOrCreate(&users[i]).Error; err != nil {
			log.Fatalf("Failed to seed user %s: %v", users[i].Email, err)
		}
	}

	log.Printf("Seeded %d users", len(users))
	return users
}

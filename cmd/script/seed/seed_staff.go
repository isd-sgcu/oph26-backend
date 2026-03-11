package main

import (
	"log"
	"oph26-backend/internal/entity"

	"gorm.io/gorm"
)

// seedStaff seeds 2 staff members linked to users[3] and users[4].
func seedStaff(db *gorm.DB, users []entity.User) []entity.Staff {
	staffRecords := []entity.Staff{
		{
			UserID:    &users[3].ID,
			Cuid:      "STF001",
			Firstname: "สมชาย",
			Surname:   "ใจดี",
			Nickname:  "ชาย",
			Phone:     "0812345678",
			Year:      "4",
			Email:     users[3].Email,
			Faculty:   "eng",
		},
		{
			UserID:    &users[4].ID,
			Cuid:      "STF002",
			Firstname: "สมหญิง",
			Surname:   "รักไทย",
			Nickname:  "หญิง",
			Phone:     "0898765432",
			Year:      "3",
			Email:     users[4].Email,
			Faculty:   "sci",
		},
	}

	for i := range staffRecords {
		if err := db.Where("cuid = ?", staffRecords[i].Cuid).FirstOrCreate(&staffRecords[i]).Error; err != nil {
			log.Fatalf("Failed to seed staff %s: %v", staffRecords[i].Cuid, err)
		}
	}

	log.Printf("Seeded %d staff", len(staffRecords))
	return staffRecords
}

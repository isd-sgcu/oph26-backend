package main

import (
	"log"
	"oph26-backend/internal/entity"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// seedAttendees seeds 3 attendees linked to users[0..2].
// attendees[0] and [1] are highschool (eligible for MyPiece).
// attendees[2] is a parent.
func seedAttendees(db *gorm.DB, users []entity.User, staffs []entity.Staff) []entity.Attendee {
	checkinTime := time.Now()

	attendees := []entity.Attendee{
		{
			UserID:                        users[0].ID,
			Firstname:                     "นภา",
			Surname:                       "ดาวเรือง",
			AttendeeType:                  "highschool",
			Age:                           17,
			Province:                      "กรุงเทพมหานคร",
			StudyLevel:                    strPtr("มัธยมศึกษาปีที่ 5"),
			SchoolName:                    strPtr("โรงเรียนสาธิตจุฬาฯ"),
			NewsSourceSelected:            pq.StringArray{"facebook", "line"},
			InitialFirstInterestedFaculty: "eng",
			InterestedFaculty:             pq.StringArray{"eng", "sci"},
			ObjectiveSelected:             pq.StringArray{"explore"},
			TicketCode:                    "H000001",
			CheckinStaffID:                &staffs[0].ID,
			CheckedInAt:                   &checkinTime,
		},
		{
			UserID:                        users[1].ID,
			Firstname:                     "ธนา",
			Surname:                       "วิชาการ",
			AttendeeType:                  "highschool",
			Age:                           16,
			Province:                      "เชียงใหม่",
			StudyLevel:                    strPtr("มัธยมศึกษาปีที่ 4"),
			SchoolName:                    strPtr("โรงเรียนยุพราชวิทยาลัย"),
			NewsSourceSelected:            pq.StringArray{"instagram"},
			InitialFirstInterestedFaculty: "md",
			InterestedFaculty:             pq.StringArray{"md", "pharm"},
			ObjectiveSelected:             pq.StringArray{"explore"},
			TicketCode:                    "H000002",
		},
		{
			UserID:                        users[2].ID,
			Firstname:                     "มานี",
			Surname:                       "กุลดี",
			AttendeeType:                  "parent",
			Age:                           45,
			Province:                      "ขอนแก่น",
			InitialFirstInterestedFaculty: "law",
			InterestedFaculty:             pq.StringArray{"law", "arts"},
			ObjectiveSelected:             pq.StringArray{"explore"},
			TicketCode:                    "P000001",
		},
	}

	for i := range attendees {
		if err := db.Where("ticket_code = ?", attendees[i].TicketCode).FirstOrCreate(&attendees[i]).Error; err != nil {
			log.Fatalf("Failed to seed attendee %s: %v", attendees[i].TicketCode, err)
		}
	}

	log.Printf("Seeded %d attendees", len(attendees))
	return attendees
}

func strPtr(s string) *string {
	return &s
}

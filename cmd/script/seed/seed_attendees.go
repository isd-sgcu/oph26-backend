package main

import (
	"log"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// seedAttendees seeds 3 attendees linked to users[0..2].
// attendees[0] and [1] are highschool (eligible for MyPiece).
// attendees[2] is a parent.
func seedAttendees(db *gorm.DB, users []entity.User, staffs []entity.Staff) []entity.Attendee {
	matthayomPlai := model.MatthayomPlai
	attendees := []entity.Attendee{
		{
			UserID:                        users[0].ID,
			Firstname:                     "นภา",
			Surname:                       "ดาวเรือง",
			AttendeeType:                  "student",
			DateOfBirth:                   time.Date(2009, 3, 17, 0, 0, 0, 0, time.UTC),
			Province:                      "กรุงเทพมหานคร",
			StudyLevel:                    &matthayomPlai,
			SchoolName:                    strPtr("โรงเรียนสาธิตจุฬาฯ"),
			NewsSourceSelected:            pq.StringArray{"facebook", "line"},
			InitialFirstInterestedFaculty: model.ENG,
			InterestedFaculty:             entity.FacultyList{model.ENG, model.SCI},
			ObjectiveSelected:             pq.StringArray{"explore"},
			TicketCode:                    "H000001",
		},
		{
			UserID:                        users[1].ID,
			Firstname:                     "ธนา",
			Surname:                       "วิชาการ",
			AttendeeType:                  "student",
			DateOfBirth:                   time.Date(2010, 3, 17, 0, 0, 0, 0, time.UTC),
			Province:                      "เชียงใหม่",
			StudyLevel:                    &matthayomPlai,
			SchoolName:                    strPtr("โรงเรียนยุพราชวิทยาลัย"),
			NewsSourceSelected:            pq.StringArray{"instagram"},
			InitialFirstInterestedFaculty: model.MD,
			InterestedFaculty:             entity.FacultyList{model.MD, model.PHARM},
			ObjectiveSelected:             pq.StringArray{"explore"},
			TicketCode:                    "H000002",
		},
		{
			UserID:                        users[2].ID,
			Firstname:                     "มานี",
			Surname:                       "กุลดี",
			AttendeeType:                  "parent",
			DateOfBirth:                   time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC),
			Province:                      "ขอนแก่น",
			InitialFirstInterestedFaculty: model.LAW,
			InterestedFaculty:             entity.FacultyList{model.LAW, model.ARTS},
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

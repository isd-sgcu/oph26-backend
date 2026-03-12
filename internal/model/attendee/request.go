package attendee

import "github.com/lib/pq"

type AttendeeCreateRequest struct {
	Age                  int            `json:"age" validate:"min=5,max=100"`
	AttendeeType         string         `json:"attendee_type" validate:"oneof=student parent educationstaff other"`
	Firstname            string         `json:"firstname" validate:"max=200"`
	InterestedFaculty    pq.StringArray `json:"interested_faculty" validate:"min=1,max=4,unique"`
	NewsSourcesOther     *string        `json:"news_sources_other" validate:"omitempty,max=200"`
	NewsSourceSelected   pq.StringArray `json:"news_sources_selected" validate:"min=1"`
	ObjectiveOther       *string        `json:"objective_other" validate:"omitempty,max=200"`
	ObjectiveSelected    pq.StringArray `json:"objective_selected" validate:"min=1"`
	Province             string         `json:"province" validate:"max=200"`
	SchoolName           *string        `json:"school_name" validate:"omitempty,max=200"`
	StudyLevel           *string        `json:"study_level" validate:"omitempty,max=200"`
	Surname              string         `json:"surname" validate:"max=200"`
	TransportationMethod string         `json:"transportation_method" validate:"max=200"` // api docs said "ไม่ต้อง validation ก็ได้"
}

type PutAttendeesRequest struct {
	Firstname            *string         `json:"firstname" validate:"omitempty,max=200"`
	Surname              *string         `json:"surname" validate:"omitempty,max=200"`
	Age                  *int            `json:"age" validate:"omitempty,min=5,max=100"`
	Province             *string         `json:"province" validate:"omitempty,max=200"`
	StudyLevel           *string         `json:"study_level" validate:"omitempty,max=200"`
	SchoolName           *string         `json:"school_name" validate:"omitempty,max=200"`
	NewsSourceSelected   *pq.StringArray `json:"news_sources_selected" validate:"omitempty,min=1"`
	NewsSourcesOther     *string         `json:"news_sources_other" validate:"omitempty,max=200"`
	InterestedFaculty    *pq.StringArray `json:"interested_faculty" validate:"omitempty,min=1,max=4,unique"`
	ObjectiveSelected    *pq.StringArray `json:"objective_selected" validate:"omitempty,min=1"`
	ObjectiveOther       *string         `json:"objective_other" validate:"omitempty,max=200"`
	TransportationMethod *string         `json:"transportation_method" validate:"omitempty, max=200"` // api docs said "ไม่ต้อง validation ก็ได้"
}

type PutFavoriteWorkshopsRequest struct {
	Code  string `json:"code"`
	State bool   `json:"state"`
}

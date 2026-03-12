package attendee

type AttendeeCreateRequest struct {
	DateOfBirth          string    `json:"date_of_birth"`
	AttendeeType         string    `json:"attendee_type" validate:"oneof=student parent educationstaff other"`
	Firstname            string    `json:"firstname" validate:"max=200"`
	InterestedFaculty    *[]string `json:"interested_faculty" validate:"omitempty,min=1,max=4,unique"`
	NewsSourcesOther     *string   `json:"news_sources_other" validate:"omitempty,max=200"`
	NewsSourceSelected   []string  `json:"news_sources_selected" validate:"min=1"`
	ObjectiveOther       *string   `json:"objective_other" validate:"omitempty,max=200"`
	ObjectiveSelected    []string  `json:"objective_selected" validate:"min=1"`
	Province             string    `json:"province" validate:"max=200"`
	District             string    `json:"district" validate:"max=200"`
	SchoolName           *string   `json:"school_name" validate:"omitempty,max=200"`
	StudyLevel           *string   `json:"study_level" validate:"omitempty,max=200"`
	Surname              string    `json:"surname" validate:"max=200"`
	TransportationMethod string    `json:"transportation_method" validate:"max=200"` // api docs said "ไม่ต้อง validation ก็ได้"
}

type PutAttendeesRequest struct {
	Firstname            *string   `json:"firstname" validate:"omitempty,max=200"`
	Surname              *string   `json:"surname" validate:"omitempty,max=200"`
	DateOfBirth          *string   `json:"date_of_birth" validate:"omitempty"`
	Province             *string   `json:"province" validate:"omitempty,max=200"`
	District             *string   `json:"district" validate:"omitempty,max=200"`
	StudyLevel           *string   `json:"study_level" validate:"omitempty,max=200"`
	SchoolName           *string   `json:"school_name" validate:"omitempty,max=200"`
	NewsSourceSelected   *[]string `json:"news_sources_selected" validate:"omitempty,min=1"`
	NewsSourcesOther     *string   `json:"news_sources_other" validate:"omitempty,max=200"`
	InterestedFaculty    *[]string `json:"interested_faculty" validate:"omitempty,min=1,max=4,unique"`
	ObjectiveSelected    *[]string `json:"objective_selected" validate:"omitempty,min=1"`
	ObjectiveOther       *string   `json:"objective_other" validate:"omitempty,max=200"`
	TransportationMethod *string   `json:"transportation_method" validate:"omitempty,max=200"` // api docs said "ไม่ต้อง validation ก็ได้"
}

type PutFavoriteWorkshopsRequest struct {
	Code  string `json:"code"`
	State bool   `json:"state"`
}

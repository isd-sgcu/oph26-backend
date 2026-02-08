package attendee

import "github.com/lib/pq"

type PutAttendeesRequest struct {
	Firstname          *string         `json:"firstname" validate:"max=200"`
	Surname            *string         `json:"surname" validate:"max=200"`
	Age                *int            `json:"age" validate:"min=5,max=100"`
	Province           *string         `json:"province" validate:"max=200"`
	StudyLevel         *string         `json:"study_level" validate:"max=200"`
	SchoolName         *string         `json:"school_name" validate:"max=200"`
	NewsSourceSelected *pq.StringArray `json:"news_sources_selected" validate:"min=1"`
	NewsSourcesOther   *string         `json:"news_sources_other" validate:"max=200"`
	InterestedFaculty  *pq.StringArray `json:"interested_faculty" validate:"min=1,max=4,unique"`
	ObjectiveSelected  *pq.StringArray `json:"objective_selected" validate:"min=1"`
	ObjectiveOther     *string         `json:"objective_other" validate:"max=200"`
}

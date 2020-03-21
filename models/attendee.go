package models

import (
    "fmt"
    "io"
    "strconv"
    "time"
)

//Attendee is used to show attribtues of attendees.
type Attendee struct {
    ID        string `json:"id"`
    UserID    string `json:"user_id"`
    MeetupID  string `json:"meetup_id"`
    Status    string `json:"status"`
    Meetup    *Meetup
    User      *User
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"-" pg:",soft_delete"`
}

//CreateAttendanceInput is used to specify attributes which are required for creating attendee.
type CreateAttendanceInput struct {
    Status   AttendanceStatus `json:"status"`
    MeetupID string           `json:"meetup_id"`
}

//AttendanceStatus type.
type AttendanceStatus string

const (
    //AttendanceStatusGoing is a type of enum status
    AttendanceStatusGoing AttendanceStatus = "going"
    // AttendanceStatusInterested is a type of enum status
    AttendanceStatusInterested AttendanceStatus = "interested"
)

// AllAttendanceStatus conatins the enum values.
var AllAttendanceStatus = []AttendanceStatus{
    AttendanceStatusGoing,
    AttendanceStatusInterested,
}

//IsValid is used to check against the value of attendance status is valid or not.
func (e AttendanceStatus) IsValid() bool {
    switch e {
    case AttendanceStatusGoing, AttendanceStatusInterested:
        return true
    }
    return false
}

func (e AttendanceStatus) String() string {
    return string(e)
}

//UnmarshalGQL is used to unmarshal the graphql values.
func (e *AttendanceStatus) UnmarshalGQL(v interface{}) error {
    str, ok := v.(string)
    if !ok {
        return fmt.Errorf("enums must be strings")
    }

    *e = AttendanceStatus(str)
    if !e.IsValid() {
        return fmt.Errorf("%s is not a valid AttendanceStatus", str)
    }
    return nil
}

//MarshalGQL is used to marshal the values to graphql.
func (e AttendanceStatus) MarshalGQL(w io.Writer) {
    fmt.Fprint(w, strconv.Quote(e.String()))
}

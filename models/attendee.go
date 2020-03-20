package models

import "time"

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

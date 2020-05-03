package models

import (
    "fmt"
    "io"
    "strconv"
    "time"
)

//Notification is used to show attribtues of attendees.
type Notification struct {
    ID             string         `json:"id"`
    UserID         string         `json:"user_id"`
    NotifiableType NotifiableType `json:"notifiable_type"`
    NotifiableID   string         `json:"notifiable_id"`
    User           *User
    ReadAt         time.Time  `json:"read_at"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
    DeletedAt      *time.Time `json:"-" pg:",soft_delete"`
}

type NotifiableType string

const (
    NotifiableTypeReplyCreated        NotifiableType = "reply_created"
    NotifiableTypeCommentCreated      NotifiableType = "comment_created"
    NotifiableTypeMeetupCreated       NotifiableType = "meetup_created"
    NotifiableTypeMeetupReminder      NotifiableType = "meetup_reminder"
    NotifiableTypeMeetupSharedToGroup NotifiableType = "meetup_shared_to_group"
)

var AllNotifiableType = []NotifiableType{
    NotifiableTypeReplyCreated,
    NotifiableTypeCommentCreated,
    NotifiableTypeMeetupCreated,
    NotifiableTypeMeetupReminder,
    NotifiableTypeMeetupSharedToGroup,
}

func (e NotifiableType) IsValid() bool {
    switch e {
    case NotifiableTypeReplyCreated, NotifiableTypeCommentCreated, NotifiableTypeMeetupCreated, NotifiableTypeMeetupReminder, NotifiableTypeMeetupSharedToGroup:
        return true
    }
    return false
}

func (e NotifiableType) String() string {
    return string(e)
}

func (e *NotifiableType) UnmarshalGQL(v interface{}) error {
    str, ok := v.(string)
    if !ok {
        return fmt.Errorf("enums must be strings")
    }

    *e = NotifiableType(str)
    if !e.IsValid() {
        return fmt.Errorf("%s is not a valid NotifiableType", str)
    }
    return nil
}

func (e NotifiableType) MarshalGQL(w io.Writer) {
    fmt.Fprint(w, strconv.Quote(e.String()))
}

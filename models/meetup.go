package models

import (
    "time"

    "github.com/go-playground/validator"
)

//Meetup model.
type Meetup struct {
    ID          string      `json:"id"`
    Name        string      `json:"name" validate:"required,min=3,max=100"`
    Description string      `json:"description" `
    UserID      string      `json:"user_id"`
    StartDate   string      `json:"start_date"`
    EndDate     string      `json:"end_date"`
    Location    string      `json:"location"`
    Categories  []*Category `pg:"many2many:category_meetup,joinFK:category_id"`
    Attendees   []*Attendee
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"-" pg:",soft_delete"`
}

//CreateMeetupInput is used to validate the attributes with the following criteria.
type CreateMeetupInput struct {
    Name        string `json:"name" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"required,min=3,max=500"`
    StartDate   string `json:"start_date" validate:"required"`
    EndDate     string `json:"end_date" validate:"required"`
    Location    string `json:"location" validate:"required,min=3,max=100"`
    GroupID     string `json:"group_id" validate:"omitempty"`
}

//UpdateMeetupInput is used to validate the attributes with the following criteria.
type UpdateMeetupInput struct {
    Name        string `json:"name" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"required,min=3,max=500"`
    StartDate   string `json:"start_date" validate:"required"`
    EndDate     string `json:"end_date" validate:"requireed"`
    Location    string `json:"location" validate:"required,min=3,max=100"`
}

//MeetupFilterInput is used to specify the attributes needed to filter meetups by.
type MeetupFilterInput struct {
    Name      *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
    StartDate *string `json:"start_date,omitempty"`
    EndDate   *string `json:"end_date,omitempty"`
    Location  *string `json:"location,omitempty" validate:"omitempty,min=3,max=100"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *UpdateMeetupInput) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

//Validate is used to validate the passed values against the struct validation props.
func (m *CreateMeetupInput) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

//Validate is used to validate the passed values of against the struct validation props.
func (m *MeetupFilterInput) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

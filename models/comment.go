package models

import (
    "time"

    "github.com/go-playground/validator"
)

//Comment model.
type Comment struct {
    ID        string     `json:"id"`
    Body      string     `json:"body"`
    UserID    string     `json:"user_id"`
    MeetupID  string     `json:"meetup_id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"-" pg:",soft_delete"`
}

//CreateCommentInput is used to validate against the attributes.
type CreateCommentInput struct {
    Body     string `json:"body" validate:"required,min=3,max=500"`
    MeetupID string `json:"meetup_id"`
}

//UpdateCommentInput is used to validate against the attributes.
type UpdateCommentInput struct {
    Body string `json:"body" validate:"required,min=3,max=500"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *CreateCommentInput) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

//Validate is used to validate the passed values against the struct validation props.
func (u *UpdateCommentInput) Validate() error {
    validate := validator.New()
    return validate.Struct(u)
}

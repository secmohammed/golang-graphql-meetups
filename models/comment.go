package models

import "github.com/go-playground/validator"

//Comment model.
type Comment struct {
    ID       string `json:"id"`
    Body     string `json:"body" validate:"required,min=3,max=500"`
    UserID   string `json:"user_id"`
    MeetupID string `json:"meetup_id"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *Comment) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

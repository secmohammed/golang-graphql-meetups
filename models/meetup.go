package models

import "github.com/go-playground/validator"

//Meetup model.
type Meetup struct {
    ID          string `json:"id"`
    Name        string `json:"name" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"required,min=3,max=500"`
    UserID      string `json:"user_id"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *Meetup) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

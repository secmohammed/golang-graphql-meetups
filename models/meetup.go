package models

import "github.com/go-playground/validator"

type Meetup struct {
    ID          string `json:"id"`
    Name        string `json:"name" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"required,min=3,max=500"`
    UserID      string `json:"userId"`
}

func (m *Meetup) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

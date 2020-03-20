package models

import (
    "time"

    "github.com/go-playground/validator"
)

//Category model attributes.
type Category struct {
    ID      string    `json:"id"`
    Name    string    `json:"name"`
    UserID  string    `json:"user_id"`
    Meetups []*Meetup `pg:"many2many:category_meetup,joinFK:meetup_id"`
    Users   []*User   `pg:"many2many:category_user,joinFK:user_id"`

    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"-" pg:",soft_delete"`
}

// CategoryMeetup struct type
type CategoryMeetup struct {
    CategoryID string
    MeetupID   string
}

//CreateCategoryInput is used to validate against the attributes.
type CreateCategoryInput struct {
    Name string `json:"name" validate:"required,min=3,max=32"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *CreateCategoryInput) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

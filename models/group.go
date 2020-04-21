package models

import (
    "time"

    "github.com/go-playground/validator"
    "github.com/secmohammed/meetups/utils/validation"
)

//Group model attributes.
type Group struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    UserID      string `json:"user_id"`
    User        *User
    Meetups     []*Meetup   `pg:"many2many:group_meetup,joinFK:meetup_id"`
    Members     []*User     `pg:"many2many:group_user,joinFK:user_id"`
    Categories  []*Category `pg:"many2many:category_group"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
    DeletedAt   *time.Time  `json:"-" pg:",soft_delete"`
}
type UserGroup struct {
    User *User  `json:"user"`
    Type string `json:"type"`
}

// GroupUser struct type
type GroupUser struct {
    tableName struct{} `sql:"group_user"`
    UserID    string
    GroupID   string
    Type      string
}

// CategoryGroup struct type
type CategoryGroup struct {
    tableName struct{} `sql:"category_group"`

    CategoryID string
    GroupID    string
}
type MeetupGroup struct {
    tableName struct{} `sql:"group_meetup"`
    GroupID   string
    MeetupID  string
}

//UpdateGroupInput is used to validate against the attributes.
type UpdateGroupInput struct {
    Name        string   `json:"name" validate:"required,min=3,max=32"`
    Description string   `json:"description" validate:"required,min=3,max=32"`
    CategoryIds []string `json:"category_ids" validate:"omitempty,is_slice,is_string_element"`
}

//CreateGroupInput is used to validate against the attributes.
type CreateGroupInput struct {
    Name        string   `json:"name" validate:"omitempty,min=3,max=32"`
    Description string   `json:"description" validate:"omitempty,min=3,max=32"`
    CategoryIds []string `json:"category_ids" validate:"omitempty,is_slice,is_string_element"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *CreateGroupInput) Validate() error {
    validate := validator.New()
    validate.RegisterValidation("is_slice", validation.IsSlice)
    validate.RegisterValidation("is_string_element", validation.IsStringElem)
    return validate.Struct(m)
}

//Validate is used to validate the passed values against the struct validation props.
func (m *UpdateGroupInput) Validate() error {
    validate := validator.New()
    validate.RegisterValidation("is_slice", validation.IsSlice)
    validate.RegisterValidation("is_string_element", validation.IsStringElem)
    return validate.Struct(m)
}

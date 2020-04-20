package models

import (
    "time"

    "github.com/go-playground/validator"
)

//Group model attributes.
type Group struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    UserID      string `json:"user_id"`
    User        *User
    Meetups     []*Meetup   `pg:"many2many:category_meetup,joinFK:meetup_id"`
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
    CategoryID string
    GroupID    string
}

//UpdateGroupInput is used to validate against the attributes.
type UpdateGroupInput struct {
    Name        string   `json:"name" validate:"required,min=3,max=32"`
    Description string   `json:"description" validate:"required,min=3,max=32"`
    Categories  []string `json:"categories"`
}

//CreateGroupInput is used to validate against the attributes.
type CreateGroupInput struct {
    Name        string   `json:"name" validate:"required,min=3,max=32"`
    Description string   `json:"description" validate:"required,min=3,max=32"`
    Categories  []string `json:"categories"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *CreateGroupInput) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

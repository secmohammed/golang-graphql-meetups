package models

import (
    "fmt"
    "io"
    "strconv"
    "time"

    "github.com/go-playground/validator"
)

//Group model attributes.
type Group struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Description string     `json:"description"`
    UserID      string     `json:"user_id"`
    Meetups     []*Meetup  `pg:"many2many:category_meetup,joinFK:meetup_id"`
    Users       []*User    `pg:"many2many:category_user"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"-" pg:",soft_delete"`
}
type UserGroup struct {
    User *User  `json:"user"`
    Type string `json:"type"`
}

// GroupUser struct type
type GroupUser struct {
    UserID  string
    GroupID string
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

type Role string

const (
    RoleMember    Role = "MEMBER"
    RoleAdmin     Role = "ADMIN"
    RoleModerator Role = "MODERATOR"
)

var AllRole = []Role{
    RoleMember,
    RoleAdmin,
    RoleModerator,
}

func (e Role) IsValid() bool {
    switch e {
    case RoleMember, RoleAdmin, RoleModerator:
        return true
    }
    return false
}

func (e Role) String() string {
    return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
    str, ok := v.(string)
    if !ok {
        return fmt.Errorf("enums must be strings")
    }

    *e = Role(str)
    if !e.IsValid() {
        return fmt.Errorf("%s is not a valid Role", str)
    }
    return nil
}

func (e Role) MarshalGQL(w io.Writer) {
    fmt.Fprint(w, strconv.Quote(e.String()))
}

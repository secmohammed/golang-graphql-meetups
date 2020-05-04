package models

import "time"

type Role struct {
    ID          string
    Permissions map[string]bool
    Name        string
    Slug        string
    UserID      string  `json:"user_id"`
    Users       []*User `pg:"many2many:role_user,joinFK:user_id"`
    User        *User
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"-" pg:",soft_delete"`
}
type RoleUser struct {
    tableName struct{} `sql:"role_user"`
    RoleID    string
    UserID    string
}

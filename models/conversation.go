package models

import (
    "time"

    "github.com/go-playground/validator"
    "github.com/secmohammed/meetups/utils/validation"
)

//Conversation model.
type Conversation struct {
    ID        string `json:"id"`
    Message   string `json:"message"`
    UserID    string `json:"user_id"`
    User      *User  `json:"user"`
    ParentID  string `json:"parent_id"`
    LastReply string `json:"last_reply"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
    // Conversations []*Conversation `json:"conversations" pg:"hasmany:conversations"`
    Users     []*User    `json:"users" pg:"many2many:conversation_user"`
    DeletedAt *time.Time `json:"-" pg:",soft_delete"`
}
type CreateConversationInput struct {
    UserIds []string `json:"user_ids" validate:"required,is_slice,is_string_element"`
    Message string   `json:"message" validate:"required,min=3,max=100"`
}

//CreateMessageInput validation.
type CreateMessageInput struct {
    Message string `json:"message" validate:"required,min=3,max=100"`
}
type ConversationUser struct {
    tableName      struct{} `sql:"conversation_user"`
    UserID         string   `json:"user_id"`
    ConversationID string   `json:"conversation_id"`
}

//Validate is used to validate the passed values against the struct validation props.
func (m *CreateMessageInput) Validate() error {
    validate := validator.New()
    return validate.Struct(m)
}

//Validate is used to validate the passed values against the struct validation props.
func (m *CreateConversationInput) Validate() error {
    validate := validator.New()
    validate.RegisterValidation("is_slice", validation.IsSlice)
    validate.RegisterValidation("is_string_element", validation.IsStringElem)
    return validate.Struct(m)
}

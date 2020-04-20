package postgres

import (
    "fmt"

    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// ConversationsRepo is used to contain the db driver.
type ConversationsRepo struct {
    DB *pg.DB
}

func (c *ConversationsRepo) CreateConversationUsers(userIds []string, conversation *models.Conversation) error {
    conversations := make([]models.ConversationUser, len(userIds))
    fmt.Println(conversations, len(userIds))
    for i := 0; i < len(userIds); i++ {
        conversations = append(conversations, models.ConversationUser{
            UserID:         userIds[i],
            ConversationID: conversation.ID,
        })

    }
    _, err := c.DB.Model(&conversations).Insert()
    return err
}

//Create is used to create a conversation using the passed struct.
func (c *ConversationsRepo) Create(conversation *models.Conversation) (*models.Conversation, error) {
    _, err := c.DB.Model(conversation).Returning("*").Insert()
    return conversation, err
}

// Update is used to update the passed conversation by id.
func (c *ConversationsRepo) Update(conversation *models.Conversation) (*models.Conversation, error) {
    _, err := c.DB.Model(conversation).Where("id = ?", conversation.ID).Update()
    return conversation, err
}

// GetByID is used to fetch conversation by id.
func (c *ConversationsRepo) GetByID(id string) (*models.Conversation, error) {
    conversation := models.Conversation{}
    err := c.DB.Model(&conversation).Where("id = ?", id).Select()
    if err != nil {
        return nil, err
    }
    return &conversation, nil
}

// Delete is used to delete conversation by its id.
func (c *ConversationsRepo) Delete(conversation *models.Conversation) error {
    _, err := c.DB.Model(conversation).Where("id = ?", conversation.ID).Delete()
    return err
}

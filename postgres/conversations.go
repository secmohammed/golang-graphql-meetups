package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// ConversationsRepo is used to contain the db driver.
type ConversationsRepo struct {
    DB *pg.DB
}

func (c *ConversationsRepo) CreateConversationUsers(userIds []string, conversation *models.Conversation) error {
    conversations := make([]models.ConversationUser, len(userIds))
    for i := 0; i < len(userIds); i++ {
        conversations = append(conversations, models.ConversationUser{
            UserID:         userIds[i],
            ConversationID: conversation.ID,
        })

    }
    _, err := c.DB.Model(&conversations).Insert()
    return err
}
func (c *ConversationsRepo) IsUserParticipantOfConversation(userID, conversationID string) (bool, error) {
    conversationUser := models.ConversationUser{}
    err := c.DB.Model(&conversationUser).Where("user_id = ?", userID).Where("conversation_id = ?", conversationID).Select()
    if err != nil {
        return false, err
    }
    return true, nil

}
func (c *ConversationsRepo) GetConversationParticipants(conversationID string) ([]*models.ConversationUser, error) {
    var participants []*models.ConversationUser
    err := c.DB.Model(&participants).Where("conversation_id = ?", conversationID).Select()
    if err != nil {
        return nil, err
    }
    return participants, nil
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
func (c *ConversationsRepo) GetConversationMessages(id string) ([]*models.Conversation, error) {
    var conversations []*models.Conversation
    err := c.DB.Model(&conversations).Where("parent_id = ?", id).Select()
    if err != nil {
        return nil, err
    }
    return conversations, nil
}

// Delete is used to delete conversation by its id.
func (c *ConversationsRepo) Delete(conversation *models.Conversation) error {
    _, err := c.DB.Model(conversation).Where("id = ?", conversation.ID).Delete()
    return err
}

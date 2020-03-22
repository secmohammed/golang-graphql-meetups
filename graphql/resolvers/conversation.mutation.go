package resolvers

import (
    "context"
    "time"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

func (c *mutationResolver) CreateConversation(ctx context.Context, input models.CreateConversationInput) (*models.Conversation, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    if err := input.Validate(); err != nil {
        return nil, err
    }
    // create the first message with authenticated user.
    newConversation := &models.Conversation{
        UserID:  currentUser.ID,
        Message: input.Message,
    }
    conversation, err := c.ConversationsRepo.Create(newConversation)
    if err != nil {
        return nil, err
    }
    conversationUsers := []*models.ConversationUser{}
    users := []*models.User{}
    // prepare to insert to the conversation_user table, and prepare the user ids to find with.
    for _, id := range input.UserIds {
        conversationUsers = append(conversationUsers, &models.ConversationUser{
            UserID:         id,
            ConversationID: conversation.ID,
        })
        users = append(users, &models.User{ID: id})
    }
    //TODO:  take the ids of users and check against their existence

    // then insert them at conversation_user table
    if _, err = c.ConversationsRepo.CreateConversationUsers(conversationUsers); err != nil {
        return nil, err
    }
    // return the conversation.
    return conversation, nil
}
func (c *mutationResolver) CreateMessage(ctx context.Context, conversationID string, input models.CreateMessageInput) (*models.Conversation, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    if err := input.Validate(); err != nil {
        return nil, err
    }
    if _, err := c.ConversationsRepo.GetByID(conversationID); err != nil {
        return nil, err
    }
    newConversation := &models.Conversation{
        ParentID:  conversationID,
        UserID:    currentUser.ID,
        Message:   input.Message,
        LastReply: time.Now().Format(time.RFC3339),
    }
    conversation, err := c.ConversationsRepo.Create(newConversation)
    if err != nil {
        return nil, err
    }
    for _, observer := range conversationMessageAdded {
        observer <- conversation
    }

    return conversation, nil
}

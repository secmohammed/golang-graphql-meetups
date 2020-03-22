package resolvers

import (
    "context"
    "time"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
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
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, errors.ErrUnauthenticated
    }
    if err := input.Validate(); err != nil {
        return nil, err
    }
    _, err = c.ConversationsRepo.GetByID(conversationID)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    newConversation := &models.Conversation{
        ParentID:  conversationID,
        UserID:    currentUser.ID,
        Message:   input.Message,
        LastReply: time.Now().String(),
    }
    return c.ConversationsRepo.Create(newConversation)
}

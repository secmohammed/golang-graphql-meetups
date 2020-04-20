package resolvers

import (
    "context"
    "fmt"
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
    conversation := &models.Conversation{
        UserID:  currentUser.ID,
        Message: input.Message,
    }
    conversation, err := c.ConversationsRepo.Create(conversation)
    if err != nil {
        return nil, err
    }
    // then insert them at conversation_user table
    if err = c.ConversationsRepo.CreateConversationUsers(input.UserIds, conversation); err != nil {
        fmt.Println("here")
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
        // implement logic here to push the conversation to only people related.
        observer <- conversation
    }

    return conversation, nil
}

package resolvers

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/appleboy/go-fcm"
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
    //TODO: validate against existence of passed users first.

    // then insert them at conversation_user table
    if err = c.ConversationsRepo.CreateConversationUsers(input.UserIds, conversation); err != nil {
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
    conversation := &models.Conversation{
        ParentID:  conversationID,
        UserID:    currentUser.ID,
        Message:   input.Message,
        LastReply: time.Now().Format(time.RFC3339),
    }
    conversation, err := c.ConversationsRepo.Create(conversation)
    if err != nil {
        return nil, err
    }
    c.nClient.Publish("conversation."+conversationID, conversation)
    msg := &fcm.Message{
        To: currentUser.DeviceToken,
        Data: map[string]interface{}{
            "conversation": conversation,
            "message":      fmt.Sprintf("%s messaged you with: %s", currentUser.Username, conversation.Message),
        },
    }
    client, err := fcm.NewClient(os.Getenv("fcm_token"))
    if err != nil {
        log.Fatalln(err)
    }

    // Send the message and receive the response without retries.
    _, err = client.Send(msg)
    if err != nil {
        log.Fatalln(err)
    }
    return conversation, nil
}

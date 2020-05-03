package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

type conversationResolver struct{ *Resolver }

func (r *Resolver) Conversation() graphql.ConversationResolver {
    return &conversationResolver{r}
}
func (r *conversationResolver) Conversations(ctx context.Context, conversation *models.Conversation) ([]*models.Conversation, error) {
    return r.ConversationsRepo.GetConversationMessages(conversation.ID)
}

type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Subscription() graphql.SubscriptionResolver {
    return &subscriptionResolver{r}
}
func (r *subscriptionResolver) MessageAdded(ctx context.Context, conversationID string) (<-chan *models.Conversation, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    value, err := r.ConversationsRepo.IsUserParticipantOfConversation(currentUser.ID, conversationID)
    if err != nil || !value {
        return nil, errors.ErrUnauthenticated
    }
    event := make(chan *models.Conversation, 1)
    sub, err := r.nClient.Subscribe("conversation."+conversationID, func(t *models.Conversation) {
        event <- t
    })
    if err != nil {
        return nil, err
    }

    go func() {
        <-ctx.Done()
        sub.Unsubscribe()
    }()
    return event, nil
}
func (r *queryResolver) Conversation(ctx context.Context, id string) (*models.Conversation, error) {
    return r.ConversationsRepo.GetByID(id)
}

func (r *queryResolver) Conversations(ctx context.Context) ([]*models.Conversation, error) {
    // TODO: implement conversations
    return nil, nil
}

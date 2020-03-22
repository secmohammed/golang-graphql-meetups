package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils"
)

type conversationResolver struct{ *Resolver }

func (r *Resolver) Conversation() graphql.ConversationResolver {
    return &conversationResolver{r}
}
func (r *conversationResolver) Conversations(context.Context, *models.Conversation) ([]*models.Conversation, error) {
    return nil, nil
}

var conversationMessageAdded map[string]chan *models.Conversation

func init() {
    conversationMessageAdded = map[string]chan *models.Conversation{}
}

type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Subscription() graphql.SubscriptionResolver {
    return &subscriptionResolver{r}
}
func (r *subscriptionResolver) MessageAdded(ctx context.Context) (<-chan *models.Conversation, error) {
    id := utils.GenerateRandomString(8)

    conversationEvent := make(chan *models.Conversation, 1)
    go func() {
        <-ctx.Done()
    }()
    event := conversationEvent
    conversationMessageAdded[id] = event
    return conversationEvent, nil
}

func (r *queryResolver) Conversation(ctx context.Context, id string) (*models.Conversation, error) {
    return r.ConversationsRepo.GetByID(id)
}

func (r *queryResolver) Conversations(ctx context.Context) ([]*models.Conversation, error) {
    return nil, nil
}

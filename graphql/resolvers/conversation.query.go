package resolvers

import (
    "context"
    "log"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
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
func (r *subscriptionResolver) MessageAdded(ctx context.Context) (<-chan *models.Conversation, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    event := make(chan *models.Conversation, 1)
    sub, err := r.nClient.Subscribe("conversation", func(t *models.Conversation) {
        participants, err := r.ConversationsRepo.GetConversationParticipants(t.ParentID)

        if err != nil {
            log.Fatalln("couldn't find participants", err)
        }
        for _, participant := range participants {
            if currentUser.ID == participant.UserID {
                event <- t
            }
        }
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

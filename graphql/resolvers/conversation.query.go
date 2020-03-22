package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils"
)

type conversationResolver struct{ *Resolver }

func (r *queryResolver) Conversations(ctx context.Context) ([]*models.Conversation, error) {
    return nil, nil
}

type Chatroom struct {
    Name          string
    Conversations []models.Conversation
    Observers     map[string]struct {
        ID           string
        Conversation chan *models.Conversation
    }
}

type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Subscription() graphql.SubscriptionResolver {
    return &subscriptionResolver{r}
}
func (r *subscriptionResolver) MessageAdded(ctx context.Context, roomName string) (<-chan *models.Conversation, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    r.mu.Lock()
    room := r.Rooms[roomName]
    if room == nil {
        room = &Chatroom{
            Name: roomName,
            Observers: map[string]struct {
                ID           string
                Conversation chan *models.Conversation
            }{},
        }
        r.Rooms[roomName] = room
    }
    r.mu.Unlock()

    id := utils.GenerateRandomString(8)
    events := make(chan *models.Conversation, 1)

    go func() {
        <-ctx.Done()
        r.mu.Lock()
        delete(room.Observers, id)
        r.mu.Unlock()
    }()

    r.mu.Lock()
    room.Observers[id] = struct {
        ID           string
        Conversation chan *models.Conversation
    }{ID: currentUser.ID, Conversation: events}
    r.mu.Unlock()

    return events, nil
}

func (r *queryResolver) Conversation(ctx context.Context, id string) (*models.Conversation, error) {
    return r.ConversationsRepo.GetByID(id)
}

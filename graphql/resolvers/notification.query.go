package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
)

type notificationResolver struct{ *Resolver }

func (r *Resolver) Notification() graphql.NotificationResolver {
    return &notificationResolver{r}
}
func (n *notificationResolver) User(ctx context.Context, obj *models.Notification) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)
}
func (r *subscriptionResolver) NotificationPushed(ctx context.Context) (<-chan *models.Notification, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    event := make(chan *models.Notification, 1)
    sub, err := r.nClient.Subscribe("notification.user_"+currentUser.ID, func(t *models.Notification) {
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
func (r *queryResolver) Notifications(ctx context.Context) ([]*models.Notification, error) {
    return nil, nil
}
func (r *queryResolver) Notification(ctx context.Context, id string) (*models.Notification, error) {
    return r.NotificationsRepo.GetByID(id)
}

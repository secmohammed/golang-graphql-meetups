package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

type userResolver struct{ *Resolver }

func (r *Resolver) User() graphql.UserResolver {
    return &userResolver{r}
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
    return r.UsersRepo.GetByID(id)
}
func (u *userResolver) Roles(ctx context.Context, obj *models.User) ([]*models.Role, error) {
    return loaders.GetLoaders(ctx).RolesByUser.Load(obj.ID)
}
func (r *queryResolver) FilteredMeetupsForUser(ctx context.Context, filter *models.MeetupFilterInput, limit *int, offset *int) ([]*models.Meetup, error) {
    currentUser, err := middlewares.GetCurrentUserFromContext(ctx)
    if err != nil {
        return nil, errors.ErrUnauthenticated
    }
    if err := filter.Validate(); err != nil {
        return nil, err
    }

    return r.MeetupsRepo.GetFilteredMeetupsBasedOnUser(currentUser.ID, filter, limit, offset)
}
func (u *userResolver) Notifications(ctx context.Context, obj *models.User) ([]*models.Notification, error) {
    return loaders.GetLoaders(ctx).NotificationsByUser.Load(obj.ID)
}

func (r *queryResolver) AuthenticatedUser(ctx context.Context) (*models.User, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    return r.UsersRepo.GetByID(currentUser.ID)
}
func (u *userResolver) Permissions(ctx context.Context, obj *models.User) (interface{}, error) {
    return obj.Permissions, nil
}
func (u *userResolver) Meetups(ctx context.Context, obj *models.User) ([]*models.Meetup, error) {
    return u.MeetupsRepo.GetMeetupsForUser(obj)
}
func (u *userResolver) Interests(ctx context.Context, obj *models.User) ([]*models.Category, error) {
    return loaders.GetLoaders(ctx).InterestsByUser.Load(obj.ID)
}
func (u *userResolver) Comments(ctx context.Context, obj *models.User) ([]*models.Comment, error) {
    return u.CommentsRepo.GetCommentsForUser(obj)
}

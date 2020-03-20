package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

type userResolver struct{ *Resolver }

func (r *Resolver) User() graphql.UserResolver {
    return &userResolver{r}
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
    return r.UsersRepo.GetByID(id)
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

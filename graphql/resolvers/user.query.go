package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
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

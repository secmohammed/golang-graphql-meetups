package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql/resolvers"
    "github.com/secmohammed/meetups/models"
)

type meetupResolver struct{ *resolvers.Resolver }

func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {

    return m.Resolver.UsersRepo.GetByID(obj.UserID)
}

func (r *Resolver) Meetup() MeetupResolver {
    return &meetupResolver{r}
}

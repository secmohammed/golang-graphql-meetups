package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

func (r *Resolver) Meetup() graphql.MeetupResolver {
    return &meetupResolver{r}
}

func (r *queryResolver) Meetups(ctx context.Context, filter *graphql.MeetupFilter, limit *int, offset *int) ([]*models.Meetup, error) {

    return r.MeetupsRepo.GetMeetups(filter, limit, offset)
}

type meetupResolver struct{ *Resolver }

func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {

    return loaders.GetUserLoader(ctx).Load(obj.UserID)
}

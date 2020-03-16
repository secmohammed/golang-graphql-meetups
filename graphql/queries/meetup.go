package queries

import (
    "context"

    "github.com/secmohammed/meetups/graphql/resolvers"
    "github.com/secmohammed/meetups/models"
)

type queryResolver struct{ *resolvers.Resolver }

func (r *queryResolver) Meetups(ctx context.Context) ([]*models.Meetup, error) {

    return r.MeetupsRepo.GetMeetups()
}

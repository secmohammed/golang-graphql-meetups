package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

type attendeeResolver struct{ *Resolver }

func (r *Resolver) Attendee() graphql.AttendeeResolver {
    return &attendeeResolver{r}
}

func (a *attendeeResolver) User(ctx context.Context, obj *models.Attendee) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)
}
func (a *attendeeResolver) Meetup(ctx context.Context, obj *models.Attendee) (*models.Meetup, error) {
    return loaders.GetLoaders(ctx).MeetupByID.Load(obj.MeetupID)
}

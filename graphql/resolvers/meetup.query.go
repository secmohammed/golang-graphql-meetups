package resolvers

import (
    "context"

    "github.com/secmohammed/meetups/graphql"
    "github.com/secmohammed/meetups/graphql/loaders"
    "github.com/secmohammed/meetups/models"
)

//Meetup resolver.
func (r *Resolver) Meetup() graphql.MeetupResolver {
    return &meetupResolver{r}
}

func (r *queryResolver) Meetups(ctx context.Context, filter *models.MeetupFilter, limit *int, offset *int) ([]*models.Meetup, error) {
    return r.MeetupsRepo.GetMeetups(filter, limit, offset)
}
func (r *queryResolver) Meetup(ctx context.Context, id string) (*models.Meetup, error) {
    return r.MeetupsRepo.GetByID(id)
}

type meetupResolver struct{ *Resolver }

func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {
    return loaders.GetLoaders(ctx).UserByID.Load(obj.UserID)
}
func (m *meetupResolver) Attendees(ctx context.Context, obj *models.Meetup) ([]*models.Attendee, error) {
    return loaders.GetLoaders(ctx).AttendeesByMeetup.Load(obj.ID)
}
func (m *meetupResolver) Comments(ctx context.Context, obj *models.Meetup) ([]*models.Comment, error) {
    // loaders.GetLoaders(ctx).CommentsByMeetupID.Load(obj.ID)
    return m.CommentsRepo.GetCommentsForMeetup(obj.ID)
}
func (m *meetupResolver) Categories(ctx context.Context, obj *models.Meetup) ([]*models.Category, error) {
    return loaders.GetLoaders(ctx).CategoriesByMeetup.Load(obj.ID)
}
